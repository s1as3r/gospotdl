package download

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/bogem/id3v2"
	"github.com/s1as3r/gospotdl/search"
)

// setId3tags is a helper function that sets the id3 tags
// of a mp3 file as provided by Song `s`s
func setId3Tags(filePath string, s *search.Song) error {
	tag, err := id3v2.Open(filePath, id3v2.Options{})
	if err != nil {
		return fmt.Errorf("Error while opening file(%s): %s\n", filePath, err)
	}
	defer tag.Close()

	tag.SetTitle(s.Name)
	tag.SetArtist(s.Artists[0].Name)
	tag.SetAlbum(s.Album.Name)
	tag.SetYear(s.Album.ReleaseDate)

	resp, err := http.Get(s.Album.Images[0].URL)
	if err != nil {
		return fmt.Errorf("Error Getting Cover Image: %s\n", err)
	}
	defer resp.Body.Close()

	img, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading cover image: %s", err)
	}
	albumCover := id3v2.PictureFrame{
		Encoding:    tag.DefaultEncoding(),
		MimeType:    "image/jpeg",
		PictureType: id3v2.PTFrontCover,
		Description: "Front Cover",
		Picture:     img,
	}
	tag.AddAttachedPicture(albumCover)

	trackNumberTag := id3v2.NewEmptyTag()
	trackNumberFrame := id3v2.TextFrame{
		Encoding: id3v2.EncodingUTF8,
		Text:     strconv.Itoa(s.TrackNumber),
	}
	trackNumberTag.AddFrame("TRCK", trackNumberFrame)

	if err = tag.Save(); err != nil {
		return fmt.Errorf("Error while saving a tag: %s", err)
	}
	return nil
}
