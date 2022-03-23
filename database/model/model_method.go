package model

func (p *Product_Category) IsEqual(pc *Product_Category) bool {
	return (p.ID == pc.ID)
}
