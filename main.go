package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cmd184psu/alfredo"
)

type argStruct struct {
	//verbose mode
	verbose bool
	//show version and exit
	ver        bool
	force      bool
	debug      bool
	show       bool
	dontSave   bool
	headBucket bool
	details    S3Details
}

const VERBOSE_ENV = "VERBOSE"
const FORCE_ENV = "FORCE"
const DEBUG_ENV = "DEBUG"
const CONFIG_FILE = "~/.aws-example.json"
const credentialsFile = "~/.aws/credentials"

func parseArgs() *argStruct {
	args := new(argStruct)

	flag.BoolVar(&args.verbose, "verbose", false, "verbose move")
	flag.BoolVar(&args.force, "force", false, "force (ignore errors and just generate for gencsv and other commands)")
	flag.BoolVar(&args.ver, "ver", false, "show version")
	flag.BoolVar(&args.ver, "version", false, "show version")
	flag.BoolVar(&args.debug, "debug", false, "debug mode")
	flag.BoolVar(&args.show, "show", false, "show mode")
	flag.BoolVar(&args.dontSave, "dont-save", false, "dont save the config file")
	flag.BoolVar(&args.headBucket, "head-bucket", false, "head the bucket")
	flag.StringVar(&args.details.Bucket, "bucket", "", "Override config and use this bucket instead")
	flag.StringVar(&args.details.Region, "region", "", "Override config and use this region instead")
	flag.StringVar(&args.details.Credentials.AccessKeyID, "accessKeyId", "", "Override config and use this accessKeyId instead")
	flag.StringVar(&args.details.Credentials.SecretAccessKey, "secretKey", "", "Override config and use this secretKey instead")
	flag.StringVar(&args.details.Endpoint, "endpoint", "", "Override config and use this endpoint instead")
	flag.StringVar(&args.details.Profile, "profile", "", "Override config and use this profile instead")
	flag.Parse()

	if os.Getenv(VERBOSE_ENV) != "" || args.verbose {
		args.verbose = true
		alfredo.SetVerbose(true)
	}

	alfredo.SetEnvironment(&args.verbose, VERBOSE_ENV)
	if os.Getenv(FORCE_ENV) != "" || args.force {
		args.force = true
		alfredo.SetForce(true)
	}

	alfredo.SetEnvironment(&args.force, FORCE_ENV)

	if os.Getenv(DEBUG_ENV) != "" || args.debug {
		args.debug = true
		alfredo.SetDebug(true)
	}

	alfredo.SetEnvironment(&args.debug, DEBUG_ENV)

	return args
}

func BuildVersion() string {

	alfredo.VerbosePrintln("gitbranch=" + GitBranch)
	alfredo.VerbosePrintln("ver=" + GitVersion)
	alfredo.VerbosePrintln("time=" + GitTimestamp)

	var gb string
	if strings.EqualFold(GitBranch, "main") {
		gb = ""
	} else {
		gb = "-" + GitBranch
	}

	return fmt.Sprintf("%s%s (%s)", GitVersion, gb, GitTimestamp)
}

func MergeConfigFile(configFile string, details *S3Details) error {
	alfredo.VerbosePrintln("Merging config file " + configFile)
	var altDetails S3Details

	if err := alfredo.ReadStructFromJSONFile(alfredo.ExpandTilde(configFile), &altDetails); err != nil {
		return err
	}
	if len(details.Bucket) == 0 {
		details.Bucket = altDetails.Bucket
	}
	if len(details.Region) == 0 {
		details.Region = altDetails.Region
	}
	if len(details.Credentials.AccessKeyID) == 0 {
		details.Credentials.AccessKeyID = altDetails.Credentials.AccessKeyID
	}
	if len(details.Credentials.SecretAccessKey) == 0 {
		details.Credentials.SecretAccessKey = altDetails.Credentials.SecretAccessKey
	}
	if len(details.Endpoint) == 0 {
		details.Endpoint = altDetails.Endpoint
	}
	if len(details.Profile) == 0 {
		details.Profile = altDetails.Profile
	}

	if alfredo.FileExistsEasy(alfredo.ExpandTilde(credentialsFile)) {
		alfredo.VerbosePrintln("Credentials file exists")
	}
	return nil
}

const example_version_fmt = "AWS Alfredo Example (c) C Delezenski <cmd184psu@gmail.com> - %s\n"

func main() {
	args := parseArgs()
	if args.ver {
		fmt.Printf(example_version_fmt, BuildVersion())
		os.Exit(0)
	}

	if alfredo.FileExistsEasy(alfredo.ExpandTilde(CONFIG_FILE)) {
		if err := MergeConfigFile(CONFIG_FILE, &args.details); err != nil {
			fmt.Println("Failed to merge config file: ", err)
			os.Exit(1)
		}
	}

	if args.show {
		fmt.Println(alfredo.PrettyPrint(args.details))
	}

	if args.headBucket {
		if err := args.details.HeadBucket(); err != nil {
			fmt.Println("Failed to head bucket: ", err)
			os.Exit(1)
		}
	}

	if !args.dontSave {
		if err := alfredo.WriteStructToJSONFile(alfredo.ExpandTilde(CONFIG_FILE), args.details); err != nil {
			fmt.Println("Failed to save config file: ", err)
			os.Exit(1)
		}
	}

	fmt.Println("process complete")
}
