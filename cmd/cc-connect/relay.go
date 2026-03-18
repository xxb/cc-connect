package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type relaySendOptions struct {
	from        string
	to          string
	sessionKey  string
	message     string
	dataDir     string
	timeoutSecs *int
}

type relaySendArgError struct {
	msg       string
	showUsage bool
}

var errRelaySendHelp = errors.New("relay send help requested")

func (e *relaySendArgError) Error() string {
	return e.msg
}

func runRelay(args []string) {
	if len(args) == 0 {
		printRelayUsage()
		return
	}
	switch args[0] {
	case "send":
		runRelaySend(args[1:])
	case "--help", "-h", "help":
		printRelayUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown relay subcommand: %s\n", args[0])
		printRelayUsage()
		os.Exit(1)
	}
}

func runRelaySend(args []string) {
	opts, err := parseRelaySendArgs(args)
	if err != nil {
		if errors.Is(err, errRelaySendHelp) {
			printRelaySendUsage()
			return
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		var argErr *relaySendArgError
		if errors.As(err, &argErr) && argErr.showUsage {
			printRelaySendUsage()
		}
		os.Exit(1)
	}

	sockPath := resolveSocketPath(opts.dataDir)
	if _, err := os.Stat(sockPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: cc-connect is not running (socket not found: %s)\n", sockPath)
		os.Exit(1)
	}

	payload, err := json.Marshal(struct {
		From        string `json:"from"`
		To          string `json:"to"`
		SessionKey  string `json:"session_key"`
		Message     string `json:"message"`
		TimeoutSecs *int   `json:"timeout_secs,omitempty"`
	}{
		From:        opts.from,
		To:          opts.to,
		SessionKey:  opts.sessionKey,
		Message:     opts.message,
		TimeoutSecs: opts.timeoutSecs,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	resp, err := apiPost(sockPath, "/relay/send", payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: %s\n", strings.TrimSpace(string(body)))
		os.Exit(1)
	}

	var result struct {
		Response string `json:"response"`
	}
	json.Unmarshal(body, &result)
	fmt.Print(result.Response)
}

func parseRelaySendArgs(args []string) (relaySendOptions, error) {
	var opts relaySendOptions
	var positional []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--from", "-f":
			if i+1 < len(args) {
				i++
				opts.from = args[i]
				continue
			}
			return opts, &relaySendArgError{msg: "--from requires a value", showUsage: true}
		case "--to", "-t":
			if i+1 < len(args) {
				i++
				opts.to = args[i]
				continue
			}
			return opts, &relaySendArgError{msg: "--to requires a value", showUsage: true}
		case "--session-key", "--session", "-s":
			if i+1 < len(args) {
				i++
				opts.sessionKey = args[i]
				continue
			}
			return opts, &relaySendArgError{msg: "--session-key requires a value", showUsage: true}
		case "--message", "-m":
			if i+1 < len(args) {
				i++
				opts.message = args[i]
				continue
			}
			return opts, &relaySendArgError{msg: "--message requires a value", showUsage: true}
		case "--data-dir":
			if i+1 < len(args) {
				i++
				opts.dataDir = args[i]
				continue
			}
			return opts, &relaySendArgError{msg: "--data-dir requires a value", showUsage: true}
		case "--timeout", "--timeout-secs":
			if i+1 < len(args) {
				i++
				timeoutSecs, err := parseRelayTimeoutSecs(args[i])
				if err != nil {
					return opts, err
				}
				opts.timeoutSecs = timeoutSecs
				continue
			}
			return opts, &relaySendArgError{msg: "--timeout requires a value", showUsage: true}
		case "--help", "-h":
			return opts, errRelaySendHelp
		default:
			positional = append(positional, args[i])
		}
	}

	if opts.from == "" {
		opts.from = os.Getenv("CC_PROJECT")
	}
	if opts.sessionKey == "" {
		opts.sessionKey = os.Getenv("CC_SESSION_KEY")
	}
	if opts.message == "" && len(positional) > 0 {
		if opts.to == "" && len(positional) >= 2 {
			opts.to = positional[0]
			opts.message = strings.Join(positional[1:], " ")
		} else {
			opts.message = strings.Join(positional, " ")
		}
	}

	if opts.to == "" || opts.message == "" {
		return opts, &relaySendArgError{msg: "target project (--to) and message are required", showUsage: true}
	}
	if opts.sessionKey == "" {
		return opts, &relaySendArgError{msg: "session key is required (set CC_SESSION_KEY or use --session-key)"}
	}

	return opts, nil
}

func parseRelayTimeoutSecs(raw string) (*int, error) {
	secs, err := strconv.Atoi(raw)
	if err != nil {
		return nil, &relaySendArgError{msg: fmt.Sprintf("invalid --timeout value %q", raw), showUsage: true}
	}
	if secs < 0 {
		return nil, &relaySendArgError{msg: "timeout must be >= 0", showUsage: true}
	}
	return &secs, nil
}

func printRelayUsage() {
	fmt.Println(`Usage: cc-connect relay <command> [options]

Commands:
  send      Send a message to another bot via relay

Run 'cc-connect relay <command> --help' for details.`)
}

func printRelaySendUsage() {
	fmt.Println(`Usage: cc-connect relay send [options] [<target_project> <message>]

Send a message to another bot and wait for the response.

Options:
  -f, --from <project>       Source project (auto-detected from CC_PROJECT env)
  -t, --to <project>         Target bot project name
  -s, --session-key <key>    Session key (auto-detected from CC_SESSION_KEY env)
  -m, --message <text>       Message to send
      --timeout <seconds>    Override relay timeout for this send (0 = disable; alias: --timeout-secs)
      --data-dir <path>      Data directory (default: ~/.cc-connect)
  -h, --help                 Show this help

Examples:
  cc-connect relay send --to claude-bot "What's the weather today?"
  cc-connect relay send --timeout 300 --to gemini "Take your time"
  cc-connect relay send claude-bot What is the weather today`)
}
