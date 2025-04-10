package album

import "errors"

var ErrPlantAlreadyInAlbum = errors.New("plant already in album")
var ErrPlantNotFound = errors.New("plant not found")

var ErrAlbumNotFound = errors.New("album not found")
