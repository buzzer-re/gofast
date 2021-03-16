package cmd


import (
	"fmt"
	"runtime"
	"errors"

	"github.com/spf13/cobra"
	"github.com/aandersonl/gofast/pkg/fastHttp"
)


type Args struct {
	NumTasks int
	OutFile	string
	Url	string
}

var cmdArgs = Args{}

var gofastCmd = &cobra.Command{
	Use: "gofast",
	Short: "A HTTP downloader accelerator using concurrency",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
		  return errors.New("Requires at least a URL argument")
		}
		cmdArgs.Url = args[0]
		return nil
	  },
	Run: func (cmd *cobra.Command, args []string) {
		var response fastHttp.FastResponse = fastHttp.GetResponse(cmdArgs.Url)
		if (cmdArgs.OutFile != "") {
			response.Filename = cmdArgs.OutFile
		}
		if response.SupportConcurrent {
		 	fmt.Printf("Starting concurrent download of %s with %d tasks\n", response.Filename, cmdArgs.NumTasks)
		 	fastHttp.ConcurrentDownload(&response, cmdArgs.NumTasks)
		} else {
		 	fmt.Printf("Remote doens't support multiple connections, downloading %s as normal file\n", response.Filename)
		 	fastHttp.NormalDownload(&response)
		}
	},
}


func Execute() {
	gofastCmd.Execute()
}

func init() {
	gofastCmd.Flags().IntVarP(&cmdArgs.NumTasks, "num-tasks","n", runtime.NumCPU() * 2, "Number of tasks to download")
	gofastCmd.Flags().StringVarP(&cmdArgs.OutFile, "output","o", "", "File output name")
}