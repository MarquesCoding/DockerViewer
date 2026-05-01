package docker

type ContainerListMsg struct {
	Items []ContainerSummary
	Err   error
}

type ContainerDetailMsg struct {
	Detail ContainerDetail
	Err    error
}

type ContainerStatsMsg struct {
	Stats StatsSnapshot
	Err   error
}

type ImageListMsg struct {
	Items []ImageSummary
	Err   error
}

type ImageDetailMsg struct {
	Detail ImageDetail
	Err    error
}

type NetworkListMsg struct {
	Items []NetworkSummary
	Err   error
}

type NetworkDetailMsg struct {
	Detail NetworkDetail
	Err    error
}

type VolumeListMsg struct {
	Items []VolumeSummary
	Err   error
}

type VolumeDetailMsg struct {
	Detail VolumeDetail
	Err    error
}
