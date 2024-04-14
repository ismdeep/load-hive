package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ismdeep/load-hive/hive"
)

func UpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "up",
		RunE: func(cmd *cobra.Command, args []string) error {
			workdir, err := os.Getwd()
			if err != nil {
				return err
			}

			h, err := hive.New(workdir, "default")
			if err != nil {
				return err
			}

			if err := h.Up(); err != nil {
				return err
			}

			return nil
		},
	}
}

func StatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "status",
		RunE: func(cmd *cobra.Command, args []string) error {
			workdir, err := os.Getwd()
			if err != nil {
				return err
			}
			h, err := hive.New(workdir, "default")
			if err != nil {
				return err
			}

			h.Status()

			return nil
		},
	}
}

func DownCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "down",
		RunE: func(cmd *cobra.Command, args []string) error {
			workdir, err := os.Getwd()
			if err != nil {
				return err
			}
			h, err := hive.New(workdir, "default")
			if err != nil {
				return err
			}
			if err := h.Down(); err != nil {
				return err
			}
			return nil
		},
	}
}

func main() {
	m := cobra.Command{
		Use:           "load-hive",
		Short:         "load-hive",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	m.AddCommand(UpCommand())
	m.AddCommand(StatusCommand())
	m.AddCommand(DownCommand())

	if err := m.Execute(); err != nil {
		fmt.Println(err)
	}
}
