package services

import (
	"gonga/app/Models"

	"gorm.io/gorm"
)

func EditMedia(db *gorm.DB, ownerID uint, ownerType string, medias []Models.Media) error {
	// Get existing media for the owner
	existingMedias := []*Models.Media{}
	if err := db.Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).Find(&existingMedias).Error; err != nil {
		return err
	}

	// Create a map to store existing media URLs for efficient lookup
	existingMediaIDs := make(map[uint]bool)
	for _, media := range existingMedias {
		existingMediaIDs[media.ID] = true
	}

	// Iterate over the updated mention user IDs
	for _, attachedMedia := range medias {
		mediaID := attachedMedia.ID

		if _, exists := existingMediaIDs[mediaID]; exists {
			// Med already exists, so remove it from the map to mark it as processed
			delete(existingMediaIDs, mediaID)
		} else {
			// Fetch the media record from the database
			var media Models.Media
			if err := db.First(&media, mediaID).Error; err != nil {
				return err
			}
			// Update the owner ID of the media to the ID of the newly created post
			media.OwnerID = ownerID
			media.OwnerType = ownerType
			db.Save(&media)
		}

		// At this point, the remaining entries in existingMentionIDs are the ones that need to be removed
	}

	// Delete the media that are no longer present in the updated media URLs
	for mediaID := range existingMediaIDs {
		if err := db.Where("owner_id = ? AND owner_type = ? AND id = ?", ownerID, ownerType, mediaID).Delete(&Models.Media{}).Error; err != nil {
			return err
		}
	}

	return nil
}
