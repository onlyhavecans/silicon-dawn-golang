/*
Copyright © 2020 Amelia Aronsohn <squirrel@wearing.black>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"onlyhavecans.works/amy/silicondawn/lib"
)

var deck lib.CardDeck

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Share the tarot of the silicon-dawn",
	Long: `Before you run this command you will need to run the get command
if you don't then it will be sad.
the tarot may give you sad messages anyways though. I choose not to stop this.`,
	Run: func(cmd *cobra.Command, args []string) {
		port := viper.GetInt("port")
		addr := "0.0.0.0:" + strconv.Itoa(port)
		cardsDirectory := viper.GetString("CardsDirectory")

		if viper.GetBool("release") {
			gin.SetMode(gin.ReleaseMode)
		}

		log.Printf("Building a deck out of %s", cardsDirectory)
		var err error
		deck, err = lib.NewCardDeck(cardsDirectory)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("We have %d cards now", deck.Count())

		router := gin.Default()
		router.LoadHTMLGlob("templates/*")
		router.GET("/", index)
		router.GET("/robots.txt", robots)
		router.Static("/cards", cardsDirectory)
		err = router.Run(addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
		log.Fatal(err)
	},
}

func index(ctx *gin.Context) {
	c, err := deck.Draw()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.HTML(http.StatusOK, "index.gohtml", gin.H{
		"dir":  "cards",
		"name": c.Front(),
		"text": c.Back(),
	})
}

func robots(ctx *gin.Context) {
	ctx.String(http.StatusOK, "User-agent: *\nDisallow: /\n")
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntP("port", "p", 3200, "Port to run on")
	err := viper.BindPFlag("Port", serveCmd.Flags().Lookup("port"))
	lib.FatalIfErr("", err)

	serveCmd.Flags().BoolP("release", "r", false, "Go-Gin release mode")
	err = viper.BindPFlag("Release", serveCmd.Flags().Lookup("release"))
	lib.FatalIfErr("", err)

	viper.SetDefault("Port", 3200)
	viper.SetDefault("Release", false)
}
