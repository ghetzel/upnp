package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ghetzel/cli"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/ghetzel/upnp"
)

func main() {
	var client *upnp.UPNP

	app := cli.NewApp()
	app.Name = `upnpfriend`
	app.Usage = `An implementation of otherwise hard-to-find UPnP client utilties`
	app.Version = `0.0.1`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   `log-level, L`,
			Usage:  `Level of log output verbosity`,
			Value:  `debug`,
			EnvVar: `LOGLEVEL`,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetLevelString(c.String(`log-level`))

		if c, err := upnp.NewUPNP(); err == nil {
			client = c
		} else {
			log.Fatalf("failed to create client: %v", err)
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:      `map`,
			Usage:     `Opens one or more local ports to the ｉｎｔｅｒｎｅｔ`,
			ArgsUsage: `1234[/{tcp|udp}] ..`,
			Action: func(c *cli.Context) {
				if xip, err := client.ExternalIPAddress(); err == nil {
					log.Infof("External IP: %v", xip)

					for _, spec := range c.Args() {
						if localport, remoteport, proto, err := parsePortSpec(spec); err == nil {
							if err := client.AddPortMapping(localport, remoteport, strings.ToUpper(proto)); err == nil {
								log.Infof("mapped 127.0.0.1:%d/%s -> %v:%d/%s", localport, proto, xip, remoteport, proto)
							} else {
								log.Fatalf("failed to map %d:%d/%s - %v", localport, remoteport, proto, err)
							}
						} else {
							log.Errorf("%q: %v", spec, err)
						}
					}
				} else {
					log.Fatalf("failed to retrieve external IP: %v", err)
				}
			},
		}, {
			Name:      `unmap`,
			Usage:     `Removes previously-mapped ports`,
			ArgsUsage: `1234[/{tcp|udp}] ..`,
			Action: func(c *cli.Context) {
				for _, spec := range c.Args() {
					if localport, _, proto, err := parsePortSpec(spec); err == nil {
						if err := client.DelPortMapping(localport, strings.ToUpper(proto)); err == nil {
							log.Infof("unmapped 127.0.0.1:%d/%s", localport, proto)
						} else {
							log.Fatalf("failed to unmap %d/%s - %v", localport, proto, err)
						}
					} else {
						log.Errorf("%q: %v", spec, err)
					}
				}
			},
		},
	}

	app.Run(os.Args)
}

func parsePortSpec(spec string) (int, int, string, error) {
	p, proto := stringutil.SplitPair(spec, `/`)
	lp, rp := stringutil.SplitPair(p, `:`)

	if rp == `` {
		rp = lp
	}

	localport := int(typeutil.Int(lp))
	remoteport := int(typeutil.Int(rp))

	if proto == `` {
		proto = `tcp`
	}

	proto = strings.ToLower(proto)

	if localport <= 0 || localport > 65535 {
		return 0, 0, ``, fmt.Errorf("invalid localport: must be an integer (0, 65535]")
	}

	if remoteport <= 0 || remoteport > 65535 {
		return 0, 0, ``, fmt.Errorf("invalid remoteport: must be an integer (0, 65535]")
	}

	switch proto {
	case `tcp`, `udp`:
		return localport, remoteport, proto, nil
	default:
		return 0, 0, ``, fmt.Errorf("invalid protocol: 'tcp' or 'udp' only")
	}
}
