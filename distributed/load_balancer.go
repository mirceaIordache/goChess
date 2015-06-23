// loadBalancer project main.go
package LoadBalancer

import (
	. "github.com/mirceaIordache/goChess/common"

	"github.com/gocircuit/circuit/client"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ThreadedStatus struct {
	Status bool
	Value  int
}

func RunRemote(args []string) int {
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Args: %q", args[0:])
	c := client.Dial(os.Args[1], nil)

	anchor, res := pickServer(c)
	t := anchor.Walk([]string{"goChess", strconv.Itoa(res + 1)})
	chess, _ := filepath.Abs(os.Args[0])
	newArgs := make([]string, len(args)+1, len(args)+1)
	newArgs[0] = t.Addr()
	copy(newArgs[1:len(newArgs)], args)

	cmd := client.Cmd{
		Path: chess,
		Args: newArgs,
	}
	p, _ := t.MakeProc(cmd)
	p.Stdin().Close()
	p.Wait()
	var result []byte
	len, _ := p.Stdout().Read(result)
	t.Scrub()

	res, _ = strconv.Atoi(string(result[:len]))
	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result %d", res)
	return res
}

func pickServer(c *client.Client) (client.Anchor, int) {

	for _, r := range c.View() {
		idle, res := isIdle(r)
		if idle {
			return r, res
		}
	}
	ChessLogger.Info("Sleeping until new server is found")
	time.Sleep(time.Duration(time.Second))
	return pickServer(c)
}

func isIdle(a client.Anchor) (bool, int) {
	t := a.View()
	load, ok := t["gochess"]
	if !ok {
		return true, 0
	}

	return len(load.View()) < 5, len(load.View())
}
