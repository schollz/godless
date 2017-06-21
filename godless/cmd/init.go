// Copyright © 2017 Johnny Morrice <john@functorama.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pkg/errors"

	lib "github.com/johnny-morrice/godless"
	"github.com/johnny-morrice/godless/api"
	"github.com/johnny-morrice/godless/internal/crypto"
	"github.com/johnny-morrice/godless/log"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise your godless environment",
	Long: `Generate crypto keys and setup your config file.

	See 'godless key' for more control`,
	Run: func(cmd *cobra.Command, args []string) {
		hash := generateKey()

		fmt.Print("Use the following hash in query signature clause:\n\n\t")
		fmt.Println(string(hash))

		flushKeysToViper()
		writeViperConfig()
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}

var keyStore api.KeyStore = lib.MakeKeyStore()

func flushKeysToViper() {
	privTexts := crypto.PrivateKeysAsText(keyStore.GetAllPrivateKeys())
	pubTexts := crypto.PublicKeysAsText(keyStore.GetAllPublicKeys())

	viper.Set(__PRIVATE_KEY_CONFIG_KEY, privTexts)
	viper.Set(__PUBLIC_KEY_CONFIG_KEY, pubTexts)
}

func readKeysFromViper() {
	maybePrivTexts := viper.Get(__PRIVATE_KEY_CONFIG_KEY)
	maybePubTexts := viper.Get(__PUBLIC_KEY_CONFIG_KEY)

	if maybePrivTexts == nil || maybePubTexts == nil {
		return
	}

	privTexts, privOk := maybePrivTexts.(string)
	pubTexts, pubOk := maybePubTexts.(string)

	if !privOk && pubOk {
		err := errors.New("Corrupt viper config for public/private keys")
		die(err)
	}

	privKeys := crypto.PrivateKeysFromText(privTexts)
	pubKeys := crypto.PublicKeysFromText(pubTexts)

	for _, pub := range pubKeys {
		keyStore.PutPublicKey(pub)
	}

	for _, priv := range privKeys {
		keyStore.PutPrivateKey(priv)
	}
}

func writeViperConfig() {
	configFilePath := viper.ConfigFileUsed()

	if configFilePath != "" {
		configFilePath = homeConfigFilePath()
	}

	file, err := os.Create(configFilePath)

	if err != nil {
		die(err)
	}

	defer func() {
		file.Close()
		err := os.Chmod(configFilePath, 0400)
		if err != nil {
			log.Error("Failed to chmod 400 key file")
		}
	}()

	writeJson(file, viper.AllSettings())
}

func writeJson(file *os.File, contents interface{}) {
	bs, err := json.MarshalIndent(contents, "", "  ")

	if err != nil {
		log.Error("Encoding JSON config failed: %v", err.Error())
	}

	w := bufio.NewWriter(file)
	_, err = w.Write(bs)

	if err != nil {
		log.Error("Error writing JSON config: %v", err.Error())
	}
}

func homeConfigFilePath() string {
	dir := os.Getenv("HOME")
	name := __CONFIG_FILE_NAME
	return path.Join(dir, name)
}

const __PRIVATE_KEY_CONFIG_KEY = "PrivateKeys"
const __PUBLIC_KEY_CONFIG_KEY = "PublicKeys"
