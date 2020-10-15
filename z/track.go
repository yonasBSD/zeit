package z

import (
  "os"
  "fmt"
  "github.com/spf13/cobra"
)

var trackCmd = &cobra.Command{
  Use:   "track",
  Short: "Tracking time",
  Long: "Add a new tracking entry, which can either be kept running until 'finish' is being called or parameterized to be a finished entry.",
  Run: func(cmd *cobra.Command, args []string) {
    user := GetCurrentUser()

    runningEntryId, err := database.GetRunningEntryId(user)
    if err != nil {
      fmt.Printf("△ %+v\n", err)
      os.Exit(1)
    }

    if runningEntryId != "" {
      fmt.Printf("▷ a task is already running\n")
      os.Exit(1)
    }

    newEntry, err := NewEntry("", begin, finish, project, task, user)
    if err != nil {
      fmt.Printf("△ %+v\n", err)
      os.Exit(1)
    }

    isRunning := newEntry.Finish.IsZero()

    _, err = database.AddEntry(user, newEntry, isRunning)
    if err != nil {
      fmt.Printf("△ %+v\n", err)
      os.Exit(1)
    }

    fmt.Printf(newEntry.GetOutputForTrack(isRunning))
    return
  },
}

func init() {
  rootCmd.AddCommand(trackCmd)
  trackCmd.Flags().StringVarP(&begin, "begin", "b", "", "Time the entry should begin at\n\nEither in the formats 16:00 / 4:00PM \nor relative to the current time, \ne.g. -0:15 (now minus 15 minutes), +1.50 (now plus 1:30h).")
  trackCmd.Flags().StringVarP(&finish, "finish", "s", "", "Time the entry should finish at\n\nEither in the formats 16:00 / 4:00PM \nor relative to the current time, \ne.g. -0:15 (now minus 15 minutes), +1.50 (now plus 1:30h).\nMust be after --begin time.")
  trackCmd.Flags().StringVarP(&project, "project", "p", "", "Project to be assigned")
  trackCmd.Flags().StringVarP(&task, "task", "t", "", "Task to be assigned")
  trackCmd.Flags().BoolVarP(&force, "force", "f", false, "Force begin tracking of a new task \neven though another one is still running \n(ONLY IF YOU KNOW WHAT YOU'RE DOING!)")

  var err error
  database, err = InitDatabase()
  if err != nil {
    fmt.Printf("△ %+v\n", err)
    os.Exit(1)
  }
}
