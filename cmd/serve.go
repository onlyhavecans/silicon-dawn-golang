package cmd

import (
	"fmt"
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
		run()
	},
}

func run() {
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
}

func index(ctx *gin.Context) {
	c, err := deck.Draw()
	if err != nil {
		err := fmt.Errorf("could not draw card: %w", err)
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
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
