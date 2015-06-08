package keypaircommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/olekukonko/tablewriter"
	osKeypairs "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var list = cli.Command{
	Name:        "list",
	Usage:       fmt.Sprintf("%s %s list [flags]", util.Name, commandPrefix),
	Description: "Lists keypairs",
	Action:      commandList,
	Flags:       util.CommandFlags(flagsList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{}
}

func commandList(c *cli.Context) {
	util.CheckArgNum(c, 0)
	client := auth.NewClient("compute")
	allPages, err := keypairs.List(client).AllPages()
	if err != nil {
		fmt.Printf("Error listing keypairs: %s\n", err)
		os.Exit(1)
	}
	o, err := osKeypairs.ExtractKeyPairs(allPages)
	if err != nil {
		fmt.Printf("Error listing keypairs: %s\n", err)
		os.Exit(1)
	}
	output.Print(c, o, tableList)
}

func tableList(c *cli.Context, i interface{}) {
	kps, ok := i.([]osKeypairs.KeyPair)
	if !ok {
		fmt.Fprintf(c.App.Writer, "Could not type assert interface\n%+v\nto []osKeypairs.KeyPair\n", i)
		os.Exit(1)
	}
	t := tablewriter.NewWriter(c.App.Writer)
	keys := []string{"Name", "Fingerprint"}
	t.SetHeader(keys)
	for _, kp := range kps {
		m := structs.Map(kp)
		f := []string{}
		for _, key := range keys {
			f = append(f, fmt.Sprint(m[key]))
		}
		t.Append(f)
	}
	t.Render()
}
