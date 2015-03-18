// loadBalancer project main.go
package LoadBalancer

import (
	"github.com/gocircuit/circuit/client"
	"os"
	"strconv"
	"time"
)

func load_balancer() {

	// The first argument is the circuit server address that this execution will use.
	// Alternative: use DialDiscover to get a *Client
	c := client.Dial(os.Args[1], nil)

	// TODO: Set up game tree with initial board setting

	// Fire-off payload processes.
	n := 1

	// TODO: Replace with ply depth reached
	for true {
		//TODO: Set correct payload and process
		cmd := client.Cmd{
			//Path: "/opt/gochess/bin/boardgenerate"
			Path: "/bin/sleep",
			Args: []string{strconv.Itoa(3 + n*3), c.Addr()},
		}
		i_ := n
		//Separate thread
		go func() {
			// Pick a circuit server to run payload on.
			t := pickServer(c).Walk([]string{"gochess", strconv.Itoa(i_)})
			// Execute the process and store it in the anchor.
			p, _ := t.MakeProc(cmd)
			// Close the process standard input to indicte no intention to write data.
			p.Stdin().Close()
			// Block until the process exits.
			p.Wait()
			// Remove the anchor storing the process.
			//var result []byte
			//p.Stdout().Read(result)
			//TODO: Add result to game tree
			t.Scrub()
			println("Payload", i_+1, "finished.")
		}()
		n++
	}
	//TODO: Run tree evaluation with recommended board to be returned
	for true {
		//TODO: Set correct payload and process
		cmd := client.Cmd{
			//Path: "/opt/gochess/bin/boardeval"
			Path: "/bin/sleep",
			Args: []string{strconv.Itoa(3 + n*3), c.Addr()},
		}
		i_ := n
		//Separate thread
		go func() {
			// Pick a circuit server to run payload on.
			t := pickServer(c).Walk([]string{"gochess", strconv.Itoa(i_)})
			// Execute the process and store it in the anchor.
			p, _ := t.MakeProc(cmd)
			// Close the process standard input to indicte no intention to write data.
			p.Stdin().Close()
			// Block until the process exits.
			p.Wait()
			// Remove the anchor storing the process.
			var result []byte
			p.Stdout().Read(result)
			//TODO: Add result to game tree
			t.Scrub()
			println("Payload", i_+1, "finished.")
		}()
		n++
	}
	println("All done.")
}

func pickServer(c *client.Client) client.Anchor {
	for _, r := range c.View() {
		if isIdle(r) {
			println(r.Addr())
			return r
		}
	}
	println("Sleeping")
	time.Sleep(time.Duration(time.Second))
	return pickServer(c)
}

func isIdle(a client.Anchor) bool {
	t := a.View()
	res, ok := t["gochess"]
	return !ok || len(res.View()) < 5
}
