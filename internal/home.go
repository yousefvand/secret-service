package internal

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// SetupHomeDirectories creates secret-service
// necessary directories under: '~/secret-service'
// .secret-service
// ├── secretservice-cli
// │   └── logs
// └── secretserviced
//     ├── config.yaml
//     ├── db.json
//     └── logs
//         └── secretserviced.log
func SetupHomeDirectories() string {

	const logHome string = "logs"
	const appHome string = ".secret-service"
	const cliHome string = "secretservice-cli"
	const serviceHome string = "secretserviced"

	// get user home directory aka '~' or global variable $HOME
	userHome, err := os.UserHomeDir()

	if err != nil {
		log.Panicf("Cannot get user home. Error: %v", err)
	}

	////////////////////////////// .secret-service //////////////////////////////

	appHomeFullPAth := filepath.Join(userHome, appHome)
	// Do we have "~/.secret-service" directory?
	appHomeExists, err := fileOrFolderExists(appHomeFullPAth)

	if err != nil {
		log.Panicf("Cannot find if app home directory exist or not: '%s'. Error : %v",
			appHomeFullPAth, err)
	}

	if !appHomeExists {
		err := os.Mkdir(appHomeFullPAth, 0755)
		if err != nil {
			log.Panicf("Cannot make '%s' directory. Error : %v", appHomeFullPAth, err)
		}
	}

	////////////////////////////// secretservice-cli //////////////////////////////

	cliHomeFullPAth := filepath.Join(appHomeFullPAth, cliHome)
	// Do we have "~/.secret-service/secretservice-cli" directory?
	cliHomeExists, err := fileOrFolderExists(cliHomeFullPAth)

	if err != nil {
		log.Panicf("Cannot find if service home directory exist or not: '%s'. Error : %v",
			cliHomeFullPAth, err)
	}

	if !cliHomeExists {
		err := os.Mkdir(cliHomeFullPAth, 0700)
		if err != nil {
			log.Panicf("Cannot make '%s' directory. Error : %v", cliHomeFullPAth, err)
		}
	}

	////////////////////////////// cli logs //////////////////////////////

	cliLogHomeFullPAth := filepath.Join(cliHomeFullPAth, logHome)
	// Do we have "~/.secret-service/secretservice-cli/logs" directory?
	cliLogHomeExists, err := fileOrFolderExists(cliLogHomeFullPAth)

	if err != nil {
		log.Panicf("Cannot find if log home directory exist or not: '%s'. Error : %v",
			cliLogHomeFullPAth, err)
	}

	if !cliLogHomeExists {
		err := os.Mkdir(cliLogHomeFullPAth, 0700)
		if err != nil {
			log.Panicf("Cannot make '%s' directory. Error : %v", cliLogHomeFullPAth, err)
		}
	}

	////////////////////////////// secretserviced //////////////////////////////

	serviceHomeFullPAth := filepath.Join(appHomeFullPAth, serviceHome)
	// Do we have "~/.secret-service/secretserviced" directory?
	serviceHomeExists, err := fileOrFolderExists(serviceHomeFullPAth)

	if err != nil {
		log.Panicf("Cannot find if service home directory exist or not: '%s'. Error : %v",
			serviceHomeFullPAth, err)
	}

	if !serviceHomeExists {
		err := os.Mkdir(serviceHomeFullPAth, 0700)
		if err != nil {
			log.Panicf("Cannot make '%s' directory. Error : %v", serviceHomeFullPAth, err)
		}
	}

	////////////////////////////// service log //////////////////////////////

	serviceLogHomeFullPAth := filepath.Join(serviceHomeFullPAth, logHome)
	// Do we have "~/.secret-service/secretserviced/logs" directory?
	serviceLogHomeExists, err := fileOrFolderExists(serviceLogHomeFullPAth)

	if err != nil {
		log.Panicf("Cannot find if log home directory exist or not: '%s'. Error : %v",
			serviceLogHomeFullPAth, err)
	}

	if !serviceLogHomeExists {
		err := os.Mkdir(serviceLogHomeFullPAth, 0700)
		if err != nil {
			log.Panicf("Cannot make '%s' directory. Error : %v", serviceLogHomeFullPAth, err)
		}
	}

	return serviceHomeFullPAth
}
