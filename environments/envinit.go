package environments

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"if0/common"
	"if0/config"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// This function checks if the environment directory contains necessary files, if not, creates them.
func envInit(envPath string) error {
	err := downloadLogo(envPath)
	if err != nil {
		fmt.Println("Error: Downloading logo -", err)
	}
	createZeroFile(envPath)
	createCIFile(envPath)
	createDash1Env(envPath)
	sshDir := filepath.Join(envPath, ".ssh")
	files, direrr := ioutil.ReadDir(sshDir)
	// .ssh dir not present or present but no keys
	if _, err := os.Stat(sshDir); os.IsNotExist(err) || (direrr == nil && len(files) < 2) {
		fmt.Printf("Creating dir %s\n", sshDir)
		_ = os.Mkdir(sshDir, 0700)
		err := generateSSHKeyPair(sshDir)
		if err != nil {
			fmt.Println("Error: Generating SSH Key pair - ", err)
			return err
		}
	}
	return nil
}

func createZeroFile(envPath string) {
	f := createFile(filepath.Join(envPath, "zero.env"))
	defer f.Close()
	if f != nil {
		repoName := strings.Replace(envPath, common.EnvDir+string(os.PathSeparator), "", 1)
		_, _ = f.WriteString("IF0_ENVIRONMENT=" + repoName + "\n")
		pwd := generateRandSeq()
		hash, err := generateHashCmd(pwd)
		if runtime.GOOS == "windows" || hash == "" || err != nil {
			hash, err = generateHashDocker(pwd)
			if err != nil {
				fmt.Println("Error: Could not create htpasswd hash -", err)
				return
			}
		}
		_, _ = f.WriteString("ZERO_ADMIN_USER=admin\n")
		_, _ = f.WriteString("ZERO_ADMIN_PASSWORD=" + pwd + "\n")
		_, _ = f.WriteString("ZERO_ADMIN_PASSWORD_HASH=" + hash + "\n")
	}
}

func createCIFile(envPath string) {
	f := createFile(filepath.Join(envPath, ".gitlab-ci.yml"))
	defer f.Close()
	if f != nil {
		shipmateUrl := getShipmateUrl()
		dataToWrite := fmt.Sprintf("include:\n  - remote: '%s'", shipmateUrl)
		_, _ = f.Write([]byte(dataToWrite))
	}
}

func createDash1Env(envPath string) {
	dash1File := createFile(filepath.Join(envPath, "dash1.env"))
	defer dash1File.Close()

	// if dash1.env is already present, do not prompt the user for dash1.env content
	var skipDash bool
	if dash1File == nil {
		skipDash = true
	}

	var dash1Content []string
	var zeroContent []string

	fmt.Print("Use Cloud Provider? [Y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	useProvider, _ := reader.ReadString('\n')
	useProvider = strings.TrimSpace(strings.ToLower(useProvider))
	if (useProvider == "y" || useProvider == "") && !skipDash {
		dash1Content = useCloudProvider(dash1Content)
	} else {
		fmt.Print("Enter IPs:")
		ips, _ := reader.ReadString('\n')
		zeroContent = append(zeroContent, "ZERO_NODES_MANAGER="+ips+"\n")
	}

	fmt.Print("Custom Domain? [y/N]: ")
	customDomain, _ := reader.ReadString('\n')
	switch strings.TrimSpace(strings.ToLower(customDomain)) {
	case "y":
		fmt.Print("Enter Domain: ")
		domain, _ := reader.ReadString('\n')
		zeroContent = append(zeroContent, "ZERO_BASE_DOMAIN="+domain+"\n")
	case "n", "":
		fmt.Println("Exiting.")
	}

	if !skipDash && len(dash1Content) > 0 {
		for _, line := range dash1Content {
			_, _ = dash1File.WriteString(line)
		}
	}

	if len(zeroContent) > 0 {
		zeroFile, _ := os.OpenFile(filepath.Join(envPath, "zero.env"), os.O_APPEND|os.O_RDWR, 0644)
		defer zeroFile.Close()
		for _, line := range zeroContent {
			_, _ = zeroFile.WriteString(line)
		}
	}
}

func useCloudProvider(dash1Content []string) []string {
	fmt.Print("Cloud Provider [1. HCLOUD, 2. Digital Ocean, 3. AWS]: ")
	reader := bufio.NewReader(os.Stdin)
	provider, _ := reader.ReadString('\n')
	switch strings.TrimSpace(strings.ToLower(provider)) {
	case "1", "":
		dash1Content = append(dash1Content, "DASH1_MODULE=hcloud\n")
		fmt.Print("Enter HCLOUD_TOKEN: ")
		htoken, _ := reader.ReadString('\n')
		dash1Content = append(dash1Content, "HCLOUD_TOKEN="+htoken)
	case "2":
		dash1Content = append(dash1Content, "DASH1_MODULE=digitalocean\n")
		fmt.Print("Enter DO_TOKEN: ")
		doToken, _ := reader.ReadString('\n')
		dash1Content = append(dash1Content, "DO_TOKEN="+doToken)
	case "3":
		dash1Content = append(dash1Content, "DASH1_MODULE=aws\n")
		fmt.Print("Enter AWS_SECRET_KEY_ID: ")
		awsSecretKeyId, _ := reader.ReadString('\n')
		dash1Content = append(dash1Content, "AWS_SECRET_KEY_ID="+awsSecretKeyId)
		fmt.Print("Enter AWS_SECRET_ACCESS_KEY: ")
		awsSecretAccessKey, _ := reader.ReadString('\n')
		dash1Content = append(dash1Content, "AWS_SECRET_ACCESS_KEY="+awsSecretAccessKey)
	}
	return dash1Content
}

func pushInitChanges(r *git.Repository, auth transport.AuthMethod) error {
	w, _ := syncObj.GetWorktree(r)
	status, _ := syncObj.Status(w)
	if len(status) > 0 {
		fmt.Println("Syncing environment init file changes")
		for file := range status {
			_ = syncObj.AddFile(w, file)
		}
		// git commit
		err := syncObj.Commit(w)
		if err != nil {
			fmt.Println("Error: Committing changes - ", err)
			return err
		}
		// git push
		err = syncObj.Push(auth, r)
		if err != nil {
			fmt.Println("Error: Pushing changes - ", err)
			return err
		}
	}
	return nil
}

func createFile(fileName string) *os.File {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Println("Creating file", fileName)
		f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
		return f
	}
	return nil
}

func getShipmateUrl() string {
	// get SHIPMATE_WORKFLOW_URL from if0.env
	config.ReadConfigFile(common.If0Default)
	shipmateUrl := config.GetEnvVariable("SHIPMATE_WORKFLOW_URL")
	// if not found, add it to if0.env and return the value
	if shipmateUrl == "" {
		f, _ := os.OpenFile(common.If0Default, os.O_APPEND, 0644)
		defer f.Close()
		_, _ = f.WriteString("SHIPMATE_WORKFLOW_URL=https://gitlab.com/peter.saarland/shipmate/-/raw/master/shipmate.gitlab-ci.yml\n")
	}
	config.ReadConfigFile(common.If0Default)
	return config.GetEnvVariable("SHIPMATE_WORKFLOW_URL")
}

func downloadLogo(envDir string) error {
	logoFile := filepath.Join(envDir, "logo.png")
	url := "https://gitlab.com/peter.saarland/scratch/-/raw/master/logo.png?inline=false"

	out, err := os.Create(logoFile)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
