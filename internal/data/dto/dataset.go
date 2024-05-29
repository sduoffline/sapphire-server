package dto

type NewDataset struct {
	Name        string `json:"name" binding:"required"`
	CreatorID   int    `json:"creator_id" binding:"required"`
	TypeID      int    `json:"type_id" binding:"required"`
	Description string `json:"description"`
	Format      string `json:"format"`
	Size        int    `json:"size"`
	IsPublic    bool   `json:"is_public"`
	IsDeleted   bool   `json:"is_deleted"`
}
