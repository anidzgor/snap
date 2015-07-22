package plugin

import (
	"encoding/gob"
	"fmt"
	"io"
	"log" // TODO proper logging to file or elsewhere
	"net"
	"net/http"
	"net/rpc"

	"github.com/intelsdi-x/pulse/control/plugin/cpolicy"
	"github.com/intelsdi-x/pulse/core/ctypes"
)

// Acts as a proxy for RPC calls to a CollectorPlugin. This helps keep the function signature simple
// within plugins vs. having to match required RPC patterns.

// Collector plugin
type CollectorPlugin interface {
	Plugin
	CollectMetrics([]PluginMetricType) ([]PluginMetricType, error)
	GetMetricTypes() ([]PluginMetricType, error)
	GetConfigPolicyTree() (cpolicy.ConfigPolicyTree, error)
}

// Execution method for a Collector plugin. Error and exit code (int) returned.
func StartCollector(c CollectorPlugin, s Session, r *Response) (error, int) {
	var exitCode int = 0

	l, err := net.Listen("tcp", "127.0.0.1:"+s.ListenPort())
	if err != nil {
		s.Logger().Println(err.Error())
		panic(err)
	}
	s.SetListenAddress(l.Addr().String())
	s.Logger().Printf("Listening %s\n", l.Addr())
	s.Logger().Printf("Session token %s\n", s.Token())

	// Create our proxy
	proxy := &collectorPluginProxy{
		Plugin:  c,
		Session: s,
	}

	// Register the proxy under the "Collector" namespace
	rpc.RegisterName("Collector", proxy)
	// Register common plugin methods used for utility reasons
	e := rpc.Register(s)
	if e != nil {
		if e.Error() != "rpc: service already defined: SessionState" {
			log.Println(e.Error())
			s.Logger().Println(e.Error())
			return e, 2
		}
	}

	switch r.Meta.RPCType {
	case JSONRPC:
		rpc.HandleHTTP()
		http.HandleFunc("/rpc", func(w http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()
			w.Header().Set("Content-Type", "application/json")
			res := NewRPCRequest(req.Body).Call()
			io.Copy(w, res)
		})
		go http.Serve(l, nil)
	case NativeRPC:
		go func() {
			for {
				conn, err := l.Accept()
				if err != nil {
					panic(err)
				}
				go rpc.ServeConn(conn)
			}
		}()
	default:
		panic("Unsupported RPC type")
	}

	resp := s.generateResponse(r)
	// Output response to stdout
	fmt.Println(string(resp))

	go s.heartbeatWatch(s.KillChan())

	if s.isDaemon() {
		exitCode = <-s.KillChan() // Closing of channel kills
	}

	return nil, exitCode
}

func init() {
	gob.Register(*(&ctypes.ConfigValueInt{}))
	gob.Register(*(&ctypes.ConfigValueStr{}))
	gob.Register(*(&ctypes.ConfigValueFloat{}))

	gob.Register(cpolicy.NewPolicyNode())
	gob.Register(&cpolicy.StringRule{})
	gob.Register(&cpolicy.IntRule{})
	gob.Register(&cpolicy.FloatRule{})
}
