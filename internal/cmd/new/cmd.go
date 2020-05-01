package new

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/chirino/graphql-gw/internal/cmd/root"
	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{
		Use:   "new",
		Short: "creates a new project with default config",
		Run:   run,
		Args:  cobra.ExactArgs(1),
	}
	ConfigFile = ""
)

func init() {
	root.Command.AddCommand(Command)
}

func run(cmd *cobra.Command, args []string) {
	dir := args[0]
	os.MkdirAll(dir, 0755)

	configFile := filepath.Join(dir, "graphql-gw.yaml")
	err := ioutil.WriteFile(configFile, []byte(`#
# Configure the host and port the service will listen on
listen: localhost:8080

#
# Configure the GraphQL upstream servers you will be accessing
upstreams:
  anilist:
    url: https://graphql.anilist.co/
    prefix: Ani

types:
  - name: Query
    actions:
      # mounts all the fields of the root anilist query onto the Query type
      - type: mount
        upstream: anilist
        query: query {}

      # mounts on a new ani_query field the root anilist query
      - type: mount
        name: ani_query
        upstream: anilist
        query: query {}

      # Adds a animeCharacters($page:Int, $perPage:Int, $search:String) field
      - type: mount
        name: animeCharacters
        upstream: anilist
        query: |
          query ($page:Int, $perPage:Int, $search:String) {
            Page(page:$page, perPage:$perPage) {
              characters(search:$search)
            }
          }

  - name: Mutation
    actions:
      # mounts all the fields of the root anilist mutation onto the Mutation type
      - type: mount
        upstream: anilist
        query: mutation {}

`), 0644)

	if err != nil {
		log.Fatalf("%+v", err)
	}

	log.Printf(`Project created in the '%s' directory.`, dir)
	log.Printf(`Edit '%s' and then run:`, configFile)
	log.Println()
	log.Println(`    cd`, dir)
	log.Println(`    graphql-gw serve`)
	log.Println()
}
