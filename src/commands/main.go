package commands

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ko6bxl/cm2img"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "test",
		Description: "Does a test function",
	},
	{
		Name:        "cm2img",
		Description: "Converts image to cm2 save string",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "upload",
				Description: "Image to Convert",
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Required:    true,
			},
		},
	},
}

var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"test": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This has been the test function, bye bye!",
			},
		})
	},
	"cm2img": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
		attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

		err := downloadFile(attachmentUrl, attachmentID)

		fileBinary, _ := os.ReadFile(attachmentID)

		mimeType := http.DetectContentType(fileBinary)

		if mimeType == "image/jpeg" {
			os.Rename(attachmentID, attachmentID+".jpeg")
			attachmentID = attachmentID + ".jpeg"
		} else if mimeType == "image/png" {
			os.Rename(attachmentID, attachmentID+".png")
			attachmentID = attachmentID + ".png"

		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to detect file type",
				},
			})
			log.Println(err)
			return
		}

		log.Println(attachmentID)
		log.Println(attachmentUrl)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to Download File",
				},
			})
			log.Println(err)
			return
		}

		out, err := cm2img.Gen("fine", "/home/me1on/Project/cm2bot/cli/"+attachmentID)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Generated File:",
				Files: []*discordgo.File{
					{
						ContentType: "text/plain",
						Name:        "save.txt",
						Reader:      strings.NewReader(out),
					},
				},
			},
		})
	},
}

func downloadFile(url, filename string) error {
	out, err := os.Create(filename)

	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return err
	}
	return nil
}
