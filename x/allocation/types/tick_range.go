package types

func (tr *TickRange) Equals(other TickRange) bool {
	return (tr.Upper == other.Upper) &&
		(tr.Lower == other.Lower) &&
		(tr.Weight == other.Weight)
}
