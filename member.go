package groupmestatsbot

import (
	"math"
	"sort"
)

// Member is a container for a GroupMe member's staistics.
type Member struct {
	ID   string
	Name string

	PopularityScore   int // how often did others favorite their messages
	UnpopularityScore int // how often did their messages get zero favorites
	SimpScore         int // how many times did they favorite someone else
	NarcissistScore   int // how many times did they favorite themselves
	VisionaryScore    int // how many images did they send
	WordsmithScore    int // how many text-only messages did they send

	Kicked   []*Member // members they've kicked from the group
	KickedBy []*Member // members who kicked them from the group

	Added   []*Member // members they've added to the group
	AddedBy []*Member // members who added them to the group

	NumMessages int // how many messages did they send
}

// NewMember creates a new Member.
func NewMember(ID, name string) *Member {
	return &Member{
		ID:       ID,
		Name:     name,
		Kicked:   []*Member{},
		KickedBy: []*Member{},
		Added:    []*Member{},
		AddedBy:  []*Member{},
	}
}

// Charisma returns the quality of their overall messages.
func (m *Member) Charisma() float64 {
	if m.PopularityScore < 1 || m.NumMessages < 1 {
		return -1
	}

	return float64(m.PopularityScore) / float64(m.NumMessages)
}

// Lurkiness returns a ratio between their interactions with others' messages and how often they post messages themselves.
func (m *Member) Lurkiness() float64 {
	if m.SimpScore < 1 || m.NumMessages < 1 {
		return -1
	}

	return float64(m.SimpScore) / float64(m.NumMessages)
}

func (s *Stats) addMember(userID, name string) {
	if m, ok := s.Members[userID]; !ok {
		s.Members[userID] = NewMember(userID, name)
	} else {
		if m.Name == "" {
			m.Name = name
		}
	}
}

func (s *Stats) incNumMessages(userID, name string) {
	s.addMember(userID, name)

	s.Members[userID].NumMessages++
}

func (s *Stats) incPopularity(userID, name string) {
	s.addMember(userID, name)

	s.Members[userID].PopularityScore++
}

func (s *Stats) incUnpopularity(userID, name string) {
	s.addMember(userID, name)

	s.Members[userID].UnpopularityScore++
}

func (s *Stats) incSimp(userID, name string) {
	s.addMember(userID, name)

	s.Members[userID].SimpScore++
}

func (s *Stats) incNarcissist(userID, name string) {
	s.addMember(userID, name)

	s.Members[userID].NarcissistScore++
}

func (s *Stats) incVisionary(userID, name string) {
	s.addMember(userID, name)

	s.Members[userID].VisionaryScore++
}

func (s *Stats) incWordsmith(userID, name string) {
	s.addMember(userID, name)

	s.Members[userID].WordsmithScore++
}

// remover/removed are temporary empty shells; they are not the true members stored in Stats.Members
func (s *Stats) addKicked(remover, removed *Member) {
	s.addMember(remover.ID, remover.Name)
	remover = s.Members[remover.ID]

	s.addMember(removed.ID, removed.Name)
	removed = s.Members[removed.ID]

	remover.Kicked = append(remover.Kicked, removed)
}

// remover/removed are temporary empty shells; they are not the true members stored in Stats.Members
func (s *Stats) addKickedBy(remover, removed *Member) {
	s.addMember(remover.ID, remover.Name)
	remover = s.Members[remover.ID]

	s.addMember(removed.ID, removed.Name)
	removed = s.Members[removed.ID]

	removed.KickedBy = append(removed.KickedBy, remover)
}

// adder/added are temporary empty shells; they are not the true members stored in Stats.Members
func (s *Stats) addAdded(adder *Member, addedUsers []*Member) {
	s.addMember(adder.ID, adder.Name)
	adder = s.Members[adder.ID]

	for _, member := range addedUsers {
		s.addMember(member.ID, member.Name)
		added := s.Members[member.ID]

		adder.Added = append(adder.Added, added)
	}
}

// adder/added are temporary empty shells; they are not the true members stored in Stats.Members
func (s *Stats) addAddedBy(adder *Member, addedUsers []*Member) {
	s.addMember(adder.ID, adder.Name)
	adder = s.Members[adder.ID]

	for _, member := range addedUsers {
		s.addMember(member.ID, member.Name)
		added := s.Members[member.ID]

		added.AddedBy = append(added.AddedBy, adder)
	}
}

// TopOfThePops returns a sorted list of the most popular members.
// Popularity is defined as someone who has the most favorites.
func (s *Stats) TopOfThePops(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].PopularityScore > sorted[j].PopularityScore })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// TopOfTheSimps returns a sorted list of the biggest simp members.
// A simp is defined as someone who favorites other members' messages the most.
func (s *Stats) TopOfTheSimps(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].SimpScore > sorted[j].SimpScore })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// TopOfTheNarcissists returns a sorted list of the biggest narcissistic members.
// A narcissist is defined as someone who favorites their own messages the most.
func (s *Stats) TopOfTheNarcissists(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].NarcissistScore > sorted[j].NarcissistScore })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// TopPosters returns a sorted list of who posted the most messages.
func (s *Stats) TopPosters(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].NumMessages > sorted[j].NumMessages })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// MostCharismatic returns a sorted list of who posts the highest quality messages.
// Charisma is defined as (# of favorites received / # of messages they posted).
func (s *Stats) MostCharismatic(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Charisma() > sorted[j].Charisma() })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// TopLurker returns a sorted list of who lurks the most.
// A lurker is defined as (# of favorites given out / # of messages they posted).
func (s *Stats) TopLurker(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Lurkiness() > sorted[j].Lurkiness() })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// TopRambler returns a sorted list of who has the most messages with zero favorites.
func (s *Stats) TopRambler(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].UnpopularityScore > sorted[j].UnpopularityScore })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// MostVisionary returns a sorted list of who posted the most images.
func (s *Stats) MostVisionary(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].VisionaryScore > sorted[j].VisionaryScore })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// TopWordsmith returns a sorted list of who posted the most text-only messages.
func (s *Stats) TopWordsmith(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].WordsmithScore > sorted[j].WordsmithScore })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// BiggestFoot returns a sorted list of who kicked the most members from the group.
func (s *Stats) BiggestFoot(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return len(sorted[i].Kicked) > len(sorted[j].Kicked) })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// SorestBum returns a sorted list of who got kicked from the group the most.
func (s *Stats) SorestBum(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return len(sorted[i].KickedBy) > len(sorted[j].KickedBy) })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// TopMother returns a sorted list of who birthed (added) the most members into the world (group).
func (s *Stats) TopMother(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return len(sorted[i].Added) > len(sorted[j].Added) })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}

// MostReincarnated returns a sorted list of who got added to the group the most.
func (s *Stats) MostReincarnated(limit int) []*Member {
	if limit == -1 {
		limit = math.MaxInt64
	}

	sorted := []*Member{}

	for _, member := range s.Members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool { return len(sorted[i].AddedBy) > len(sorted[j].AddedBy) })

	top := []*Member{}
	for i := 0; i < limit && i < len(sorted); i++ {
		top = append(top, sorted[i])
	}

	return top
}
