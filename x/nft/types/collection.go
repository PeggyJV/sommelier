package types

// NewCollection creates a new NFT Collection
func NewCollection(denom Denom, nfts ...NFT) Collection {
	collection := Collection{
		Denom: denom,
	}

	for _, nft := range nfts {
		collection = collection.AddNFT(nft.(BaseNFT))
	}

	return collection
}

// AddNFT adds an NFT to the collection
func (c Collection) AddNFT(nft BaseNFT) Collection {
	c.NFTs = append(c.NFTs, nft)
	return c
}

func (c Collection) Supply() int {
	return len(c.NFTs)
}
