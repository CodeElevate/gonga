package services

import (
	"gonga/app/Models"

	"gorm.io/gorm"
)

func EditMentions(db *gorm.DB,ownerID uint, ownerType string, mentions []Models.Mention) error {
	// Get existing mentions for the owner
	existingMentions := []Models.Mention{}
	if err := db.Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).Find(&existingMentions).Error; err != nil {
		return err
	}

	// Create a map to store existing mention user IDs for efficient lookup
	existingMentionIDs := make(map[uint]bool)
	for _, mention := range existingMentions {
		existingMentionIDs[mention.UserID] = true
	}

	// Iterate over the updated mention user IDs
	for _, mentionedUser := range mentions {
		userID := mentionedUser.UserID

		if _, exists := existingMentionIDs[userID]; exists {
			// Mention already exists, so remove it from the map to mark it as processed
			delete(existingMentionIDs, userID)
		} else {
			// Mention doesn't exist, so create a new mention
			mention := &Models.Mention{
				UserID:    userID,
				OwnerID:   ownerID,
				OwnerType: ownerType,
			}
			// Save the mention to the database
			if err := db.Create(&mention).Error; err != nil {
				return err
			}
		}

		// At this point, the remaining entries in existingMentionIDs are the ones that need to be removed
	}

	// Delete the mentions that are no longer present in the updated mentions
	for userID := range existingMentionIDs {
		if err := db.Where("owner_id = ? AND owner_type = ? AND user_id = ?", ownerID, ownerType, userID).Delete(&Models.Mention{}).Error; err != nil {
			return err
		}
	}

	return nil
}
