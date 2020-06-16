package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	Host    string
	Port    string
	Mode    string
	TimeOut int
	OutFile string
)

func init() {
	rootCmd := &cobra.Command{}
	// 程序名称
	_, f := filepath.Split(os.Args[0])
	rootCmd.Use = f
	rootCmd.Short = "ServerScan for Port Scaner and Service Version Detection"
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if err := checkArgs(); err != nil {
			cmd.HelpFunc()(cmd, args)
			log.Panic(err)
		}
	}

	rootCmd.PersistentFlags().BoolP("help", "", false, "help for this command")
	rootCmd.PersistentFlags().StringVarP(&Host, "host", "h", "", "Host to be scanned, supports four formats:\n192.168.1.1\n192.168.1.1-10\n192.168.1.*\n192.168.1.0/24.")
	rootCmd.PersistentFlags().StringVarP(&Port, "ports", "p", "80-99,7000-9000,9001-9999,4430,1433,1521,3306,5000,5432,6379,21,22,100-500,873,4440,6082,3389,5560,5900-5909,1080,1900,10809,50030,50050,50070", "Customize port list, separate with ',' example: 21,22,80-99,8000-8080 ...")
	rootCmd.PersistentFlags().StringVarP(&Mode, "mode", "m", "icmp", "Scan Mode icmp or tcp.")
	rootCmd.PersistentFlags().IntVarP(&TimeOut, "timeout", "t", 2, "Setting scaner connection timeouts,Maxtime 30 Second.")
	rootCmd.PersistentFlags().StringVarP(&OutFile, "outfile", "o", "", "Output the scanning information to file.")

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}

// 检查入参
func checkArgs() error {
	// check host
	hostsPattern := `^(([01]?\d?\d|2[0-4]\d|25[0-5])\.){3}([01]?\d?\d|2[0-4]\d|25[0-5])\/(\d{1}|[0-2]{1}\d{1}|3[0-2])$|^(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[0-9]{1,2})(\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[0-9]{1,2})){3}$`
	hostsRegexp := regexp.MustCompile(hostsPattern)
	checkHost := hostsRegexp.MatchString(Host)

	hostsPattern2 := `\b(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})\-((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2}))\b`
	hostsRegexp2 := regexp.MustCompile(hostsPattern2)
	checkHost2 := hostsRegexp2.MatchString(Host)

	hostsPattern3 := `((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(\*$)`
	hostsRegexp3 := regexp.MustCompile(hostsPattern3)
	checkHost3 := hostsRegexp3.MatchString(Host)

	if Host == "" || (checkHost == false && checkHost2 == false && checkHost3 == false) {
		return errors.New("the host is invalid")
	}

	// check port
	portsPattern := `^([0-9]|[1-9]\d|[1-9]\d{2}|[1-9]\d{3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$|^\d+(-\d+)?(,\d+(-\d+)?)*$`
	portsRegexp := regexp.MustCompile(portsPattern)
	checkPort := portsRegexp.MatchString(Port)
	if Port != "" && checkPort == false {
		return errors.New("the port is invalid")
	}

	// check mode
	if Mode != "tcp" && Mode != "icmp" {
		return errors.New("the mode is invalid")
	}

	// check timeout
	if TimeOut <= 0 || TimeOut > 30 {
		return errors.New("the timeout is invalid")
	}

	// check outFile
	if OutFile != "" && pathCheck(OutFile) == false {
		return errors.New("the outFile is invalid")
	}

	return nil
}

// 检查文件路径
func pathCheck(files string) bool {
	path, _ := filepath.Split(files)
	_, err := os.Stat(path)
	if err == nil {
		_, err2 := os.Stat(files)
		if err2 == nil {
			return false
		}
		if os.IsNotExist(err2) {
			return true
		}
	} else {
		err3 := os.MkdirAll(path, os.ModePerm)
		if err3 == nil {
			return true
		} else {
			return false
		}
	}
	return false
}
