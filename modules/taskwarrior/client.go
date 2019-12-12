package taskwarrior

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Represents a single taskwarrior instance.
type TaskWarrior struct {
	Config *TaskRC // Configuration options
	Tasks  []Task  // Task JSON entries
}

// Create new empty TaskWarrior instance.
func NewTaskWarrior(configPath string) (*TaskWarrior, error) {
	// Read the configuration file.
	taskRC, err := ParseTaskRC(configPath)
	if err != nil {
		return nil, err
	}

	// Create new TaskWarrior instance.
	tw := &TaskWarrior{Config: taskRC}
	return tw, nil
}

// Fetch all tasks for given TaskWarrior with system `taskwarrior` command call.
func (tw *TaskWarrior) FetchAllTasks() error {
	if tw == nil {
		return fmt.Errorf("Uninitialized taskwarrior database!")
	}

	rcOpt := "rc:" + tw.Config.ConfigPath
	out, err := exec.Command("task", rcOpt, "export").Output()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(out), &tw.Tasks)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
