package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var re = regexp.MustCompile(`(?m)window\[\"ytInitialData\"\] = (.*?);`)

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		set(scanner.Text())
	}
}

func set(uri string) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authority", "www.youtube.com")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.youtube.com/user/statacorp/playlists")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", "CONSENT=YES+ID.en+20170126-12-0; VISITOR_INFO1_LIVE=4JT476hhjMo; wide=1; YSC=KtGt2EkiGpE; LOGIN_INFO=AFmmF2swRQIhAMo10VidDwZBulvH9sPMILP5cQ1MksdtghMhLWkbmV4pAiBhvgml_NDmYbPrzO_oVrotJ0YAl-yN65DpJW5r5FqsAQ:QUQ3MjNmdzBsLWFVTG44VlZzOFFPSlJNM3NNRGk4cVZ5ZFpOQlRzcVNKbGpvU1RtcDBXSHlfbTV3dDgwR2Zjakh6R3lLeE5nakppMXNkNWk2TWxRX2tCa1BncGtka19pTkhiOXJoMkZQWnJHd2pfbUxTUlFyVE5JSlVFakdIUlEwRm9hTTQwR2YyMEtpblVsNWhfM3VrNEhUaVR6SUp2Zm9lbDlfV24xQi00UlpJZVdncEZVbWsw; SID=xAfRwRHePxoNNVNhbXJzEZpUv1lSdXXno_IqcYQKaH7qP5ObN_k2g5xfDOXS7oxxffdnpA.; __Secure-3PSID=xAfRwRHePxoNNVNhbXJzEZpUv1lSdXXno_IqcYQKaH7qP5ObZKnH9MLoPcUC68mNnIilIQ.; HSID=AosU4vJlL2mYO3XHc; SSID=AGvCPJ9K4VbBC7-9l; APISID=hMDERlr9iG6kuD8a/AC1U3_KIQ-ku-4MhH; SAPISID=pk0e9u10JAlDNBb6/A997nieEhTTn27nIv; __Secure-HSID=AosU4vJlL2mYO3XHc; __Secure-SSID=AGvCPJ9K4VbBC7-9l; __Secure-APISID=hMDERlr9iG6kuD8a/AC1U3_KIQ-ku-4MhH; __Secure-3PAPISID=pk0e9u10JAlDNBb6/A997nieEhTTn27nIv; PREF=f4=4000000&f5=30000&f6=400&al=en-GB+id&cvdm=grid; SIDCC=AJi4QfH7a0xjQW5BP-bCeQOLhZnxgsAvs8FA_bbVTmxjgNMuxlvo4FsNGKUUK7m01ZN7mIi3wg")
	
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	// Load the HTML document
	str, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	findString := re.FindStringSubmatch(string(str))
	playlist, err := UnmarshalPlaylist([]byte(findString[1]))
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile("/home/yuda/yudaprama/data/txt/stata.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	for _, tab := range playlist.Contents.TwoColumnBrowseResultsRenderer.Tabs {
		for _, content := range tab.TabRenderer.Content.SectionListRenderer.Contents {
			for _, rendererContent := range content.ItemSectionRenderer.Contents {
				for _, listRendererContent := range rendererContent.PlaylistVideoListRenderer.Contents {
					id := listRendererContent.PlaylistVideoRenderer.VideoID
					if len(id) != 0 {
						fmt.Println(id)
						if _, err1 := f.WriteString(","); err1 != nil {
							log.Panic(err1)
						}
						if _, err1 := f.WriteString(id); err1 != nil {
							log.Panic(err1)
						}
					}
				}
			}
		}
	}
}

func UnmarshalPlaylist(data []byte) (Playlist, error) {
	var r Playlist
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Playlist) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Playlist struct {
	ResponseContext  ResponseContext  `json:"responseContext"`
	Contents         Contents         `json:"contents"`
	TrackingParams   string           `json:"trackingParams"`
	Topbar           Topbar           `json:"topbar"`
	Microformat      Microformat      `json:"microformat"`
	Sidebar          Sidebar          `json:"sidebar"`
	FrameworkUpdates FrameworkUpdates `json:"frameworkUpdates"`
}

type Contents struct {
	TwoColumnBrowseResultsRenderer TwoColumnBrowseResultsRenderer `json:"twoColumnBrowseResultsRenderer"`
}

type TwoColumnBrowseResultsRenderer struct {
	Tabs []Tab `json:"tabs"`
}

type Tab struct {
	TabRenderer TabRenderer `json:"tabRenderer"`
}

type TabRenderer struct {
	Selected       bool               `json:"selected"`
	Content        TabRendererContent `json:"content"`
	TrackingParams string             `json:"trackingParams"`
}

type TabRendererContent struct {
	SectionListRenderer SectionListRenderer `json:"sectionListRenderer"`
}

type SectionListRenderer struct {
	Contents       []SectionListRendererContent `json:"contents"`
	TrackingParams string                       `json:"trackingParams"`
}

type SectionListRendererContent struct {
	ItemSectionRenderer ItemSectionRenderer `json:"itemSectionRenderer"`
}

type ItemSectionRenderer struct {
	Contents       []ItemSectionRendererContent `json:"contents"`
	TrackingParams string                       `json:"trackingParams"`
}

type ItemSectionRendererContent struct {
	PlaylistVideoListRenderer PlaylistVideoListRenderer `json:"playlistVideoListRenderer"`
}

type PlaylistVideoListRenderer struct {
	Contents          []PlaylistVideoListRendererContent `json:"contents"`
	PlaylistID        string                             `json:"playlistId"`
	IsEditable        bool                               `json:"isEditable"`
	CanReorder        bool                               `json:"canReorder"`
	TrackingParams    string                             `json:"trackingParams"`
	OnReorderEndpoint OnReorderEndpoint                  `json:"onReorderEndpoint"`
}

type PlaylistVideoListRendererContent struct {
	PlaylistVideoRenderer PlaylistVideoRenderer `json:"playlistVideoRenderer"`
}

type PlaylistVideoRenderer struct {
	VideoID            string                                  `json:"videoId"`
	Thumbnail          PlaylistVideoRendererThumbnail          `json:"thumbnail"`
	Title              LengthText                              `json:"title"`
	Index              ShowMoreText                            `json:"index"`
	ShortBylineText    ShortBylineText                         `json:"shortBylineText"`
	LengthText         LengthText                              `json:"lengthText"`
	NavigationEndpoint PlaylistVideoRendererNavigationEndpoint `json:"navigationEndpoint"`
	LengthSeconds      string                                  `json:"lengthSeconds"`
	TrackingParams     string                                  `json:"trackingParams"`
	IsPlayable         bool                                    `json:"isPlayable"`
	Menu               PlaylistVideoRendererMenu               `json:"menu"`
	IsWatched          bool                                    `json:"isWatched"`
	ThumbnailOverlays  []PlaylistVideoRendererThumbnailOverlay `json:"thumbnailOverlays"`
}

type ShowMoreText struct {
	SimpleText string `json:"simpleText"`
}

type LengthText struct {
	Accessibility SubscribeAccessibilityClass `json:"accessibility"`
	SimpleText    string                      `json:"simpleText"`
}

type SubscribeAccessibilityClass struct {
	AccessibilityData AccessibilityAccessibilityData `json:"accessibilityData"`
}

type AccessibilityAccessibilityData struct {
	Label string `json:"label"`
}

type PlaylistVideoRendererMenu struct {
	MenuRenderer PurpleMenuRenderer `json:"menuRenderer"`
}

type PurpleMenuRenderer struct {
	Items          []PurpleItem                `json:"items"`
	TrackingParams string                      `json:"trackingParams"`
	Accessibility  SubscribeAccessibilityClass `json:"accessibility"`
}

type PurpleItem struct {
	MenuServiceItemRenderer PurpleMenuServiceItemRenderer `json:"menuServiceItemRenderer"`
}

type PurpleMenuServiceItemRenderer struct {
	Text            Text                  `json:"text"`
	Icon            IconImage             `json:"icon"`
	ServiceEndpoint PurpleServiceEndpoint `json:"serviceEndpoint"`
	TrackingParams  string                `json:"trackingParams"`
	HasSeparator    *bool                 `json:"hasSeparator,omitempty"`
}

type IconImage struct {
	IconType string `json:"iconType"`
}

type PurpleServiceEndpoint struct {
	ClickTrackingParams          string                               `json:"clickTrackingParams"`
	CommandMetadata              OnCreateListCommandCommandMetadata   `json:"commandMetadata"`
	SignalServiceEndpoint        *PurpleSignalServiceEndpoint         `json:"signalServiceEndpoint,omitempty"`
	PlaylistEditEndpoint         *ServiceEndpointPlaylistEditEndpoint `json:"playlistEditEndpoint,omitempty"`
	AddToPlaylistServiceEndpoint *AddToPlaylistServiceEndpoint        `json:"addToPlaylistServiceEndpoint,omitempty"`
}

type AddToPlaylistServiceEndpoint struct {
	VideoID string `json:"videoId"`
}

type OnCreateListCommandCommandMetadata struct {
	WebCommandMetadata PurpleWebCommandMetadata `json:"webCommandMetadata"`
}

type PurpleWebCommandMetadata struct {
	URL      URL     `json:"url"`
	SendPost bool    `json:"sendPost"`
	APIURL   *string `json:"apiUrl,omitempty"`
}

type ServiceEndpointPlaylistEditEndpoint struct {
	PlaylistID string         `json:"playlistId"`
	Actions    []PurpleAction `json:"actions"`
}

type PurpleAction struct {
	AddedVideoID string `json:"addedVideoId"`
	Action       string `json:"action"`
}

type PurpleSignalServiceEndpoint struct {
	Signal  string         `json:"signal"`
	Actions []FluffyAction `json:"actions"`
}

type FluffyAction struct {
	AddToPlaylistCommand AddToPlaylistCommand `json:"addToPlaylistCommand"`
}

type AddToPlaylistCommand struct {
	OpenMiniplayer      bool                `json:"openMiniplayer"`
	OpenListPanel       bool                `json:"openListPanel"`
	VideoID             string              `json:"videoId"`
	ListType            string              `json:"listType"`
	OnCreateListCommand OnCreateListCommand `json:"onCreateListCommand"`
	VideoIDS            []string            `json:"videoIds"`
}

type OnCreateListCommand struct {
	ClickTrackingParams           string                             `json:"clickTrackingParams"`
	CommandMetadata               OnCreateListCommandCommandMetadata `json:"commandMetadata"`
	CreatePlaylistServiceEndpoint CreatePlaylistServiceEndpoint      `json:"createPlaylistServiceEndpoint"`
}

type CreatePlaylistServiceEndpoint struct {
	VideoIDS []string `json:"videoIds"`
	Hack     bool     `json:"hack"`
	Params   string   `json:"params"`
}

type Text struct {
	Runs []TextRun `json:"runs"`
}

type TextRun struct {
	Text string `json:"text"`
}

type PlaylistVideoRendererNavigationEndpoint struct {
	ClickTrackingParams string                  `json:"clickTrackingParams"`
	CommandMetadata     EndpointCommandMetadata `json:"commandMetadata"`
	WatchEndpoint       PurpleWatchEndpoint     `json:"watchEndpoint"`
}

type EndpointCommandMetadata struct {
	WebCommandMetadata FluffyWebCommandMetadata `json:"webCommandMetadata"`
}

type FluffyWebCommandMetadata struct {
	URL         string      `json:"url"`
	WebPageType WebPageType `json:"webPageType"`
	RootVe      int64       `json:"rootVe"`
}

type PurpleWatchEndpoint struct {
	VideoID          string `json:"videoId"`
	PlaylistID       string `json:"playlistId"`
	Index            int64  `json:"index"`
	StartTimeSeconds int64  `json:"startTimeSeconds"`
}

type ShortBylineText struct {
	Runs []ShortBylineTextRun `json:"runs"`
}

type ShortBylineTextRun struct {
	Text               string                               `json:"text"`
	NavigationEndpoint VideoOwnerRendererNavigationEndpoint `json:"navigationEndpoint"`
}

type VideoOwnerRendererNavigationEndpoint struct {
	ClickTrackingParams string                  `json:"clickTrackingParams"`
	CommandMetadata     EndpointCommandMetadata `json:"commandMetadata"`
	BrowseEndpoint      PurpleBrowseEndpoint    `json:"browseEndpoint"`
}

type PurpleBrowseEndpoint struct {
	BrowseID         ID               `json:"browseId"`
	CanonicalBaseURL CanonicalBaseURL `json:"canonicalBaseUrl"`
}

type PlaylistVideoRendererThumbnail struct {
	Thumbnails []ThumbnailElement `json:"thumbnails"`
}

type ThumbnailElement struct {
	URL    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

type PlaylistVideoRendererThumbnailOverlay struct {
	ThumbnailOverlayResumePlaybackRenderer *ThumbnailOverlayResumePlaybackRenderer `json:"thumbnailOverlayResumePlaybackRenderer,omitempty"`
	ThumbnailOverlayTimeStatusRenderer     *ThumbnailOverlayTimeStatusRenderer     `json:"thumbnailOverlayTimeStatusRenderer,omitempty"`
	ThumbnailOverlayNowPlayingRenderer     *ThumbnailOverlayNowPlayingRenderer     `json:"thumbnailOverlayNowPlayingRenderer,omitempty"`
}

type ThumbnailOverlayNowPlayingRenderer struct {
	Text Text `json:"text"`
}

type ThumbnailOverlayResumePlaybackRenderer struct {
	PercentDurationWatched int64 `json:"percentDurationWatched"`
}

type ThumbnailOverlayTimeStatusRenderer struct {
	Text  LengthText `json:"text"`
	Style string     `json:"style"`
}

type OnReorderEndpoint struct {
	ClickTrackingParams  string                                `json:"clickTrackingParams"`
	CommandMetadata      OnCreateListCommandCommandMetadata    `json:"commandMetadata"`
	PlaylistEditEndpoint OnReorderEndpointPlaylistEditEndpoint `json:"playlistEditEndpoint"`
}

type OnReorderEndpointPlaylistEditEndpoint struct {
	PlaylistID string            `json:"playlistId"`
	Actions    []TentacledAction `json:"actions"`
	Params     string            `json:"params"`
}

type TentacledAction struct {
	Action string `json:"action"`
}

type FrameworkUpdates struct {
	EntityBatchUpdate EntityBatchUpdate `json:"entityBatchUpdate"`
}

type EntityBatchUpdate struct {
	Mutations []Mutation `json:"mutations"`
}

type Mutation struct {
	EntityKey string  `json:"entityKey"`
	Type      string  `json:"type"`
	Payload   Payload `json:"payload"`
}

type Payload struct {
	SubscriptionStateEntity SubscriptionStateEntity `json:"subscriptionStateEntity"`
}

type SubscriptionStateEntity struct {
	Key        string `json:"key"`
	Subscribed bool   `json:"subscribed"`
}

type Microformat struct {
	MicroformatDataRenderer MicroformatDataRenderer `json:"microformatDataRenderer"`
}

type MicroformatDataRenderer struct {
	URLCanonical       string                         `json:"urlCanonical"`
	Title              string                         `json:"title"`
	Description        string                         `json:"description"`
	Thumbnail          PlaylistVideoRendererThumbnail `json:"thumbnail"`
	SiteName           string                         `json:"siteName"`
	AppName            string                         `json:"appName"`
	AndroidPackage     string                         `json:"androidPackage"`
	IosAppStoreID      string                         `json:"iosAppStoreId"`
	IosAppArguments    string                         `json:"iosAppArguments"`
	OgType             string                         `json:"ogType"`
	URLApplinksWeb     string                         `json:"urlApplinksWeb"`
	URLApplinksIos     string                         `json:"urlApplinksIos"`
	URLApplinksAndroid string                         `json:"urlApplinksAndroid"`
	URLTwitterIos      string                         `json:"urlTwitterIos"`
	URLTwitterAndroid  string                         `json:"urlTwitterAndroid"`
	TwitterCardType    string                         `json:"twitterCardType"`
	TwitterSiteHandle  string                         `json:"twitterSiteHandle"`
	SchemaDotOrgType   string                         `json:"schemaDotOrgType"`
	Noindex            bool                           `json:"noindex"`
	Unlisted           bool                           `json:"unlisted"`
	LinkAlternates     []LinkAlternate                `json:"linkAlternates"`
}

type LinkAlternate struct {
	HrefURL string `json:"hrefUrl"`
}

type ResponseContext struct {
	ServiceTrackingParams           []ServiceTrackingParam          `json:"serviceTrackingParams"`
	WebResponseContextExtensionData WebResponseContextExtensionData `json:"webResponseContextExtensionData"`
}

type ServiceTrackingParam struct {
	Service string  `json:"service"`
	Params  []Param `json:"params"`
}

type Param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type WebResponseContextExtensionData struct {
	YtConfigData YtConfigData `json:"ytConfigData"`
	HasDecorated bool         `json:"hasDecorated"`
}

type YtConfigData struct {
	Csn                   string `json:"csn"`
	VisitorData           string `json:"visitorData"`
	SessionIndex          int64  `json:"sessionIndex"`
	RootVisualElementType int64  `json:"rootVisualElementType"`
}

type Sidebar struct {
	PlaylistSidebarRenderer PlaylistSidebarRenderer `json:"playlistSidebarRenderer"`
}

type PlaylistSidebarRenderer struct {
	Items          []PlaylistSidebarRendererItem `json:"items"`
	TrackingParams string                        `json:"trackingParams"`
}

type PlaylistSidebarRendererItem struct {
	PlaylistSidebarPrimaryInfoRenderer   *PlaylistSidebarPrimaryInfoRenderer   `json:"playlistSidebarPrimaryInfoRenderer,omitempty"`
	PlaylistSidebarSecondaryInfoRenderer *PlaylistSidebarSecondaryInfoRenderer `json:"playlistSidebarSecondaryInfoRenderer,omitempty"`
}

type PlaylistSidebarPrimaryInfoRenderer struct {
	ThumbnailRenderer  ThumbnailRenderer                                    `json:"thumbnailRenderer"`
	Title              PlaylistSidebarPrimaryInfoRendererTitle              `json:"title"`
	Stats              []Stat                                               `json:"stats"`
	Menu               PlaylistSidebarPrimaryInfoRendererMenu               `json:"menu"`
	ThumbnailOverlays  []PlaylistSidebarPrimaryInfoRendererThumbnailOverlay `json:"thumbnailOverlays"`
	NavigationEndpoint PlaylistSidebarPrimaryInfoRendererNavigationEndpoint `json:"navigationEndpoint"`
	Description        Text                                                 `json:"description"`
	ShowMoreText       ShowMoreText                                         `json:"showMoreText"`
}

type PlaylistSidebarPrimaryInfoRendererMenu struct {
	MenuRenderer FluffyMenuRenderer `json:"menuRenderer"`
}

type FluffyMenuRenderer struct {
	Items           []FluffyItem                `json:"items"`
	TrackingParams  string                      `json:"trackingParams"`
	TopLevelButtons []TopLevelButton            `json:"topLevelButtons"`
	Accessibility   SubscribeAccessibilityClass `json:"accessibility"`
}

type FluffyItem struct {
	MenuServiceItemRenderer FluffyMenuServiceItemRenderer `json:"menuServiceItemRenderer"`
}

type FluffyMenuServiceItemRenderer struct {
	Text            Text                  `json:"text"`
	Icon            IconImage             `json:"icon"`
	ServiceEndpoint FluffyServiceEndpoint `json:"serviceEndpoint"`
	TrackingParams  string                `json:"trackingParams"`
}

type FluffyServiceEndpoint struct {
	ClickTrackingParams   string                      `json:"clickTrackingParams"`
	CommandMetadata       CommandCommandMetadata      `json:"commandMetadata"`
	SignalServiceEndpoint FluffySignalServiceEndpoint `json:"signalServiceEndpoint"`
}

type CommandCommandMetadata struct {
	WebCommandMetadata TentacledWebCommandMetadata `json:"webCommandMetadata"`
}

type TentacledWebCommandMetadata struct {
	URL      URL  `json:"url"`
	SendPost bool `json:"sendPost"`
}

type FluffySignalServiceEndpoint struct {
	Actions []StickyAction `json:"actions"`
}

type StickyAction struct {
	OpenPopupAction PurpleOpenPopupAction `json:"openPopupAction"`
}

type PurpleOpenPopupAction struct {
	Popup     PurplePopup `json:"popup"`
	PopupType string      `json:"popupType"`
}

type PurplePopup struct {
	ConfirmDialogRenderer PurpleConfirmDialogRenderer `json:"confirmDialogRenderer"`
}

type PurpleConfirmDialogRenderer struct {
	Title          Text               `json:"title"`
	TrackingParams string             `json:"trackingParams"`
	DialogMessages []Text             `json:"dialogMessages"`
	ConfirmButton  DismissButtonClass `json:"confirmButton"`
	CancelButton   DismissButtonClass `json:"cancelButton"`
}

type DismissButtonClass struct {
	ButtonRenderer DismissButtonButtonRenderer `json:"buttonRenderer"`
}

type DismissButtonButtonRenderer struct {
	Style           string                    `json:"style"`
	Size            string                    `json:"size"`
	IsDisabled      bool                      `json:"isDisabled"`
	Text            ShowMoreText              `json:"text"`
	TrackingParams  string                    `json:"trackingParams"`
	Command         *ButtonRendererCommand    `json:"command,omitempty"`
	ServiceEndpoint *TentacledServiceEndpoint `json:"serviceEndpoint,omitempty"`
}

type ButtonRendererCommand struct {
	ClickTrackingParams   string                       `json:"clickTrackingParams"`
	CommandMetadata       CommandCommandMetadata       `json:"commandMetadata"`
	SignalServiceEndpoint CommandSignalServiceEndpoint `json:"signalServiceEndpoint"`
}

type CommandSignalServiceEndpoint struct {
	Signal  string         `json:"signal"`
	Actions []IndigoAction `json:"actions"`
}

type IndigoAction struct {
	SignalAction Signal `json:"signalAction"`
}

type Signal struct {
	Signal string `json:"signal"`
}

type TentacledServiceEndpoint struct {
	ClickTrackingParams string                             `json:"clickTrackingParams"`
	CommandMetadata     OnCreateListCommandCommandMetadata `json:"commandMetadata"`
	FlagEndpoint        FlagEndpoint                       `json:"flagEndpoint"`
}

type FlagEndpoint struct {
	FlagAction string `json:"flagAction"`
}

type TopLevelButton struct {
	ToggleButtonRenderer *ToggleButtonRenderer         `json:"toggleButtonRenderer,omitempty"`
	ButtonRenderer       *TopLevelButtonButtonRenderer `json:"buttonRenderer,omitempty"`
}

type TopLevelButtonButtonRenderer struct {
	Style              string                            `json:"style"`
	Size               string                            `json:"size"`
	Icon               IconImage                         `json:"icon"`
	NavigationEndpoint *ButtonRendererNavigationEndpoint `json:"navigationEndpoint,omitempty"`
	Accessibility      AccessibilityAccessibilityData    `json:"accessibility"`
	Tooltip            string                            `json:"tooltip"`
	TrackingParams     string                            `json:"trackingParams"`
	IsDisabled         *bool                             `json:"isDisabled,omitempty"`
	ServiceEndpoint    *StickyServiceEndpoint            `json:"serviceEndpoint,omitempty"`
}

type ButtonRendererNavigationEndpoint struct {
	ClickTrackingParams string                  `json:"clickTrackingParams"`
	CommandMetadata     EndpointCommandMetadata `json:"commandMetadata"`
	WatchEndpoint       FluffyWatchEndpoint     `json:"watchEndpoint"`
}

type FluffyWatchEndpoint struct {
	VideoID    string `json:"videoId"`
	PlaylistID string `json:"playlistId"`
	Params     string `json:"params"`
}

type StickyServiceEndpoint struct {
	ClickTrackingParams        string                             `json:"clickTrackingParams"`
	CommandMetadata            OnCreateListCommandCommandMetadata `json:"commandMetadata"`
	ShareEntityServiceEndpoint ShareEntityServiceEndpoint         `json:"shareEntityServiceEndpoint"`
}

type ShareEntityServiceEndpoint struct {
	SerializedShareEntity string                              `json:"serializedShareEntity"`
	Commands              []ShareEntityServiceEndpointCommand `json:"commands"`
}

type ShareEntityServiceEndpointCommand struct {
	OpenPopupAction FluffyOpenPopupAction `json:"openPopupAction"`
}

type FluffyOpenPopupAction struct {
	Popup     FluffyPopup `json:"popup"`
	PopupType string      `json:"popupType"`
	BeReused  bool        `json:"beReused"`
}

type FluffyPopup struct {
	UnifiedSharePanelRenderer UnifiedSharePanelRenderer `json:"unifiedSharePanelRenderer"`
}

type UnifiedSharePanelRenderer struct {
	TrackingParams     string `json:"trackingParams"`
	ShowLoadingSpinner bool   `json:"showLoadingSpinner"`
}

type ToggleButtonRenderer struct {
	Style                    Style                          `json:"style"`
	Size                     Size                           `json:"size"`
	IsToggled                bool                           `json:"isToggled"`
	IsDisabled               bool                           `json:"isDisabled"`
	DefaultIcon              IconImage                      `json:"defaultIcon"`
	DefaultServiceEndpoint   ServiceEndpoint                `json:"defaultServiceEndpoint"`
	ToggledIcon              IconImage                      `json:"toggledIcon"`
	ToggledServiceEndpoint   ServiceEndpoint                `json:"toggledServiceEndpoint"`
	Accessibility            AccessibilityAccessibilityData `json:"accessibility"`
	TrackingParams           string                         `json:"trackingParams"`
	DefaultTooltip           string                         `json:"defaultTooltip"`
	ToggledTooltip           string                         `json:"toggledTooltip"`
	AccessibilityData        SubscribeAccessibilityClass    `json:"accessibilityData"`
	ToggledAccessibilityData SubscribeAccessibilityClass    `json:"toggledAccessibilityData"`
}

type ServiceEndpoint struct {
	ClickTrackingParams string                             `json:"clickTrackingParams"`
	CommandMetadata     OnCreateListCommandCommandMetadata `json:"commandMetadata"`
	LikeEndpoint        LikeEndpoint                       `json:"likeEndpoint"`
}

type LikeEndpoint struct {
	Status string `json:"status"`
	Target Target `json:"target"`
}

type Target struct {
	PlaylistID string `json:"playlistId"`
}

type Size struct {
	SizeType string `json:"sizeType"`
}

type Style struct {
	StyleType string `json:"styleType"`
}

type PlaylistSidebarPrimaryInfoRendererNavigationEndpoint struct {
	ClickTrackingParams string                  `json:"clickTrackingParams"`
	CommandMetadata     EndpointCommandMetadata `json:"commandMetadata"`
	WatchEndpoint       TentacledWatchEndpoint  `json:"watchEndpoint"`
}

type TentacledWatchEndpoint struct {
	VideoID    string `json:"videoId"`
	PlaylistID string `json:"playlistId"`
}

type Stat struct {
	Runs       []TextRun `json:"runs"`
	SimpleText *string   `json:"simpleText,omitempty"`
}

type PlaylistSidebarPrimaryInfoRendererThumbnailOverlay struct {
	ThumbnailOverlaySidePanelRenderer ThumbnailOverlaySidePanelRenderer `json:"thumbnailOverlaySidePanelRenderer"`
}

type ThumbnailOverlaySidePanelRenderer struct {
	Text Text      `json:"text"`
	Icon IconImage `json:"icon"`
}

type ThumbnailRenderer struct {
	PlaylistVideoThumbnailRenderer PlaylistVideoThumbnailRenderer `json:"playlistVideoThumbnailRenderer"`
}

type PlaylistVideoThumbnailRenderer struct {
	Thumbnail PlaylistVideoRendererThumbnail `json:"thumbnail"`
}

type PlaylistSidebarPrimaryInfoRendererTitle struct {
	Runs []PurpleRun `json:"runs"`
}

type PurpleRun struct {
	Text               string                                               `json:"text"`
	NavigationEndpoint PlaylistSidebarPrimaryInfoRendererNavigationEndpoint `json:"navigationEndpoint"`
}

type PlaylistSidebarSecondaryInfoRenderer struct {
	VideoOwner VideoOwner `json:"videoOwner"`
	Button     Button     `json:"button"`
}

type Button struct {
	SubscribeButtonRenderer SubscribeButtonRenderer `json:"subscribeButtonRenderer"`
}

type SubscribeButtonRenderer struct {
	ButtonText                       Text                         `json:"buttonText"`
	SubscriberCountText              ShowMoreText                 `json:"subscriberCountText"`
	Subscribed                       bool                         `json:"subscribed"`
	Enabled                          bool                         `json:"enabled"`
	Type                             string                       `json:"type"`
	ChannelID                        ID                           `json:"channelId"`
	ShowPreferences                  bool                         `json:"showPreferences"`
	SubscriberCountWithSubscribeText ShowMoreText                 `json:"subscriberCountWithSubscribeText"`
	SubscribedButtonText             Text                         `json:"subscribedButtonText"`
	UnsubscribedButtonText           Text                         `json:"unsubscribedButtonText"`
	TrackingParams                   string                       `json:"trackingParams"`
	UnsubscribeButtonText            Text                         `json:"unsubscribeButtonText"`
	LongSubscriberCountText          Text                         `json:"longSubscriberCountText"`
	ShortSubscriberCountText         ShowMoreText                 `json:"shortSubscriberCountText"`
	SubscribeAccessibility           SubscribeAccessibilityClass  `json:"subscribeAccessibility"`
	UnsubscribeAccessibility         SubscribeAccessibilityClass  `json:"unsubscribeAccessibility"`
	NotificationPreferenceButton     NotificationPreferenceButton `json:"notificationPreferenceButton"`
	SubscribedEntityKey              string                       `json:"subscribedEntityKey"`
	OnSubscribeEndpoints             []OnSubscribeEndpoint        `json:"onSubscribeEndpoints"`
	OnUnsubscribeEndpoints           []OnUnsubscribeEndpoint      `json:"onUnsubscribeEndpoints"`
}

type NotificationPreferenceButton struct {
	SubscriptionNotificationToggleButtonRenderer SubscriptionNotificationToggleButtonRenderer `json:"subscriptionNotificationToggleButtonRenderer"`
}

type SubscriptionNotificationToggleButtonRenderer struct {
	States         []StateElement                                      `json:"states"`
	CurrentStateID int64                                               `json:"currentStateId"`
	TrackingParams string                                              `json:"trackingParams"`
	Command        SubscriptionNotificationToggleButtonRendererCommand `json:"command"`
}

type SubscriptionNotificationToggleButtonRendererCommand struct {
	CommandExecutorCommand CommandExecutorCommand `json:"commandExecutorCommand"`
}

type CommandExecutorCommand struct {
	Commands []CommandExecutorCommandCommand `json:"commands"`
}

type CommandExecutorCommandCommand struct {
	OpenPopupAction TentacledOpenPopupAction `json:"openPopupAction"`
}

type TentacledOpenPopupAction struct {
	Popup     TentacledPopup `json:"popup"`
	PopupType string         `json:"popupType"`
}

type TentacledPopup struct {
	MenuPopupRenderer MenuPopupRenderer `json:"menuPopupRenderer"`
}

type MenuPopupRenderer struct {
	Items []MenuPopupRendererItem `json:"items"`
}

type MenuPopupRendererItem struct {
	MenuServiceItemRenderer TentacledMenuServiceItemRenderer `json:"menuServiceItemRenderer"`
}

type TentacledMenuServiceItemRenderer struct {
	Text            Text                  `json:"text"`
	Icon            IconImage             `json:"icon"`
	ServiceEndpoint IndigoServiceEndpoint `json:"serviceEndpoint"`
	TrackingParams  string                `json:"trackingParams"`
	IsSelected      bool                  `json:"isSelected"`
}

type IndigoServiceEndpoint struct {
	ClickTrackingParams                         string                                      `json:"clickTrackingParams"`
	CommandMetadata                             OnCreateListCommandCommandMetadata          `json:"commandMetadata"`
	ModifyChannelNotificationPreferenceEndpoint ModifyChannelNotificationPreferenceEndpoint `json:"modifyChannelNotificationPreferenceEndpoint"`
}

type ModifyChannelNotificationPreferenceEndpoint struct {
	Params string `json:"params"`
}

type StateElement struct {
	StateID     int64      `json:"stateId"`
	NextStateID int64      `json:"nextStateId"`
	State       StateState `json:"state"`
}

type StateState struct {
	ButtonRenderer StateButtonRenderer `json:"buttonRenderer"`
}

type StateButtonRenderer struct {
	Style             string                         `json:"style"`
	Size              string                         `json:"size"`
	IsDisabled        bool                           `json:"isDisabled"`
	Icon              IconImage                      `json:"icon"`
	Accessibility     AccessibilityAccessibilityData `json:"accessibility"`
	TrackingParams    string                         `json:"trackingParams"`
	AccessibilityData SubscribeAccessibilityClass    `json:"accessibilityData"`
}

type OnSubscribeEndpoint struct {
	ClickTrackingParams string                             `json:"clickTrackingParams"`
	CommandMetadata     OnCreateListCommandCommandMetadata `json:"commandMetadata"`
	SubscribeEndpoint   SubscribeEndpoint                  `json:"subscribeEndpoint"`
}

type SubscribeEndpoint struct {
	ChannelIDS []ID   `json:"channelIds"`
	Params     string `json:"params"`
}

type OnUnsubscribeEndpoint struct {
	ClickTrackingParams   string                                     `json:"clickTrackingParams"`
	CommandMetadata       CommandCommandMetadata                     `json:"commandMetadata"`
	SignalServiceEndpoint OnUnsubscribeEndpointSignalServiceEndpoint `json:"signalServiceEndpoint"`
}

type OnUnsubscribeEndpointSignalServiceEndpoint struct {
	Signal  string           `json:"signal"`
	Actions []IndecentAction `json:"actions"`
}

type IndecentAction struct {
	OpenPopupAction StickyOpenPopupAction `json:"openPopupAction"`
}

type StickyOpenPopupAction struct {
	Popup     StickyPopup `json:"popup"`
	PopupType string      `json:"popupType"`
}

type StickyPopup struct {
	ConfirmDialogRenderer FluffyConfirmDialogRenderer `json:"confirmDialogRenderer"`
}

type FluffyConfirmDialogRenderer struct {
	TrackingParams  string       `json:"trackingParams"`
	DialogMessages  []Text       `json:"dialogMessages"`
	ConfirmButton   PurpleButton `json:"confirmButton"`
	CancelButton    PurpleButton `json:"cancelButton"`
	PrimaryIsCancel bool         `json:"primaryIsCancel"`
}

type PurpleButton struct {
	ButtonRenderer PurpleButtonRenderer `json:"buttonRenderer"`
}

type PurpleButtonRenderer struct {
	Style           string                         `json:"style"`
	Size            string                         `json:"size"`
	Text            Text                           `json:"text"`
	Accessibility   AccessibilityAccessibilityData `json:"accessibility"`
	TrackingParams  string                         `json:"trackingParams"`
	ServiceEndpoint *IndecentServiceEndpoint       `json:"serviceEndpoint,omitempty"`
}

type IndecentServiceEndpoint struct {
	ClickTrackingParams string                             `json:"clickTrackingParams"`
	CommandMetadata     OnCreateListCommandCommandMetadata `json:"commandMetadata"`
	UnsubscribeEndpoint SubscribeEndpoint                  `json:"unsubscribeEndpoint"`
}

type VideoOwner struct {
	VideoOwnerRenderer VideoOwnerRenderer `json:"videoOwnerRenderer"`
}

type VideoOwnerRenderer struct {
	Thumbnail          PlaylistVideoRendererThumbnail       `json:"thumbnail"`
	Title              VideoOwnerRendererTitle              `json:"title"`
	NavigationEndpoint VideoOwnerRendererNavigationEndpoint `json:"navigationEndpoint"`
	TrackingParams     string                               `json:"trackingParams"`
}

type VideoOwnerRendererTitle struct {
	Runs []FluffyRun `json:"runs"`
}

type FluffyRun struct {
	Text               string   `json:"text"`
	NavigationEndpoint Endpoint `json:"navigationEndpoint"`
}

type Endpoint struct {
	ClickTrackingParams string                  `json:"clickTrackingParams"`
	CommandMetadata     EndpointCommandMetadata `json:"commandMetadata"`
	BrowseEndpoint      EndpointBrowseEndpoint  `json:"browseEndpoint"`
}

type EndpointBrowseEndpoint struct {
	BrowseID string `json:"browseId"`
}

type Topbar struct {
	DesktopTopbarRenderer DesktopTopbarRenderer `json:"desktopTopbarRenderer"`
}

type DesktopTopbarRenderer struct {
	Logo                     Logo                     `json:"logo"`
	Searchbox                Searchbox                `json:"searchbox"`
	TrackingParams           string                   `json:"trackingParams"`
	CountryCode              string                   `json:"countryCode"`
	TopbarButtons            []TopbarButton           `json:"topbarButtons"`
	HotkeyDialog             HotkeyDialog             `json:"hotkeyDialog"`
	BackButton               BackButtonClass          `json:"backButton"`
	ForwardButton            BackButtonClass          `json:"forwardButton"`
	A11YSkipNavigationButton A11YSkipNavigationButton `json:"a11ySkipNavigationButton"`
}

type A11YSkipNavigationButton struct {
	ButtonRenderer A11YSkipNavigationButtonButtonRenderer `json:"buttonRenderer"`
}

type A11YSkipNavigationButtonButtonRenderer struct {
	Style          string                `json:"style"`
	Size           string                `json:"size"`
	IsDisabled     bool                  `json:"isDisabled"`
	Text           Text                  `json:"text"`
	TrackingParams string                `json:"trackingParams"`
	Command        ButtonRendererCommand `json:"command"`
}

type BackButtonClass struct {
	ButtonRenderer BackButtonButtonRenderer `json:"buttonRenderer"`
}

type BackButtonButtonRenderer struct {
	TrackingParams string                `json:"trackingParams"`
	Command        ButtonRendererCommand `json:"command"`
}

type HotkeyDialog struct {
	HotkeyDialogRenderer HotkeyDialogRenderer `json:"hotkeyDialogRenderer"`
}

type HotkeyDialogRenderer struct {
	Title          Text                          `json:"title"`
	Sections       []HotkeyDialogRendererSection `json:"sections"`
	DismissButton  DismissButtonClass            `json:"dismissButton"`
	TrackingParams string                        `json:"trackingParams"`
}

type HotkeyDialogRendererSection struct {
	HotkeyDialogSectionRenderer HotkeyDialogSectionRenderer `json:"hotkeyDialogSectionRenderer"`
}

type HotkeyDialogSectionRenderer struct {
	Title   Text     `json:"title"`
	Options []Option `json:"options"`
}

type Option struct {
	HotkeyDialogSectionOptionRenderer HotkeyDialogSectionOptionRenderer `json:"hotkeyDialogSectionOptionRenderer"`
}

type HotkeyDialogSectionOptionRenderer struct {
	Label                    Text                         `json:"label"`
	Hotkey                   string                       `json:"hotkey"`
	HotkeyAccessibilityLabel *SubscribeAccessibilityClass `json:"hotkeyAccessibilityLabel,omitempty"`
}

type Logo struct {
	TopbarLogoRenderer TopbarLogoRenderer `json:"topbarLogoRenderer"`
}

type TopbarLogoRenderer struct {
	IconImage      IconImage `json:"iconImage"`
	TooltipText    Text      `json:"tooltipText"`
	Endpoint       Endpoint  `json:"endpoint"`
	TrackingParams string    `json:"trackingParams"`
}

type Searchbox struct {
	FusionSearchboxRenderer FusionSearchboxRenderer `json:"fusionSearchboxRenderer"`
}

type FusionSearchboxRenderer struct {
	Icon            IconImage                             `json:"icon"`
	PlaceholderText Text                                  `json:"placeholderText"`
	Config          Config                                `json:"config"`
	TrackingParams  string                                `json:"trackingParams"`
	SearchEndpoint  FusionSearchboxRendererSearchEndpoint `json:"searchEndpoint"`
}

type Config struct {
	WebSearchboxConfig WebSearchboxConfig `json:"webSearchboxConfig"`
}

type WebSearchboxConfig struct {
	RequestLanguage     string `json:"requestLanguage"`
	RequestDomain       string `json:"requestDomain"`
	HasOnscreenKeyboard bool   `json:"hasOnscreenKeyboard"`
	FocusSearchbox      bool   `json:"focusSearchbox"`
}

type FusionSearchboxRendererSearchEndpoint struct {
	ClickTrackingParams string                       `json:"clickTrackingParams"`
	CommandMetadata     EndpointCommandMetadata      `json:"commandMetadata"`
	SearchEndpoint      SearchEndpointSearchEndpoint `json:"searchEndpoint"`
}

type SearchEndpointSearchEndpoint struct {
	Query string `json:"query"`
}

type TopbarButton struct {
	TopbarMenuButtonRenderer         *TopbarMenuButtonRenderer         `json:"topbarMenuButtonRenderer,omitempty"`
	NotificationTopbarButtonRenderer *NotificationTopbarButtonRenderer `json:"notificationTopbarButtonRenderer,omitempty"`
}

type NotificationTopbarButtonRenderer struct {
	Icon                      IconImage                   `json:"icon"`
	MenuRequest               MenuRequest                 `json:"menuRequest"`
	Style                     string                      `json:"style"`
	TrackingParams            string                      `json:"trackingParams"`
	Accessibility             SubscribeAccessibilityClass `json:"accessibility"`
	Tooltip                   string                      `json:"tooltip"`
	UpdateUnseenCountEndpoint UpdateUnseenCountEndpoint   `json:"updateUnseenCountEndpoint"`
	NotificationCount         int64                       `json:"notificationCount"`
	HandlerDatas              []string                    `json:"handlerDatas"`
}

type MenuRequest struct {
	ClickTrackingParams   string                             `json:"clickTrackingParams"`
	CommandMetadata       OnCreateListCommandCommandMetadata `json:"commandMetadata"`
	SignalServiceEndpoint MenuRequestSignalServiceEndpoint   `json:"signalServiceEndpoint"`
}

type MenuRequestSignalServiceEndpoint struct {
	Signal  string            `json:"signal"`
	Actions []HilariousAction `json:"actions"`
}

type HilariousAction struct {
	OpenPopupAction IndigoOpenPopupAction `json:"openPopupAction"`
}

type IndigoOpenPopupAction struct {
	Popup     IndigoPopup `json:"popup"`
	PopupType string      `json:"popupType"`
	BeReused  bool        `json:"beReused"`
}

type IndigoPopup struct {
	MultiPageMenuRenderer PopupMultiPageMenuRenderer `json:"multiPageMenuRenderer"`
}

type PopupMultiPageMenuRenderer struct {
	TrackingParams     string `json:"trackingParams"`
	Style              string `json:"style"`
	ShowLoadingSpinner bool   `json:"showLoadingSpinner"`
}

type UpdateUnseenCountEndpoint struct {
	ClickTrackingParams   string                             `json:"clickTrackingParams"`
	CommandMetadata       OnCreateListCommandCommandMetadata `json:"commandMetadata"`
	SignalServiceEndpoint Signal                             `json:"signalServiceEndpoint"`
}

type TopbarMenuButtonRenderer struct {
	Icon           *IconImage                            `json:"icon,omitempty"`
	MenuRenderer   *TopbarMenuButtonRendererMenuRenderer `json:"menuRenderer,omitempty"`
	TrackingParams string                                `json:"trackingParams"`
	Accessibility  SubscribeAccessibilityClass           `json:"accessibility"`
	Tooltip        string                                `json:"tooltip"`
	Style          *string                               `json:"style,omitempty"`
	Avatar         *Avatar                               `json:"avatar,omitempty"`
	MenuRequest    *MenuRequest                          `json:"menuRequest,omitempty"`
}

type Avatar struct {
	Thumbnails                       []ThumbnailElement               `json:"thumbnails"`
	Accessibility                    SubscribeAccessibilityClass      `json:"accessibility"`
	WebThumbnailDetailsExtensionData WebThumbnailDetailsExtensionData `json:"webThumbnailDetailsExtensionData"`
}

type WebThumbnailDetailsExtensionData struct {
	ExcludeFromVpl bool `json:"excludeFromVpl"`
}

type TopbarMenuButtonRendererMenuRenderer struct {
	MultiPageMenuRenderer MenuRendererMultiPageMenuRenderer `json:"multiPageMenuRenderer"`
}

type MenuRendererMultiPageMenuRenderer struct {
	Sections       []MultiPageMenuRendererSection `json:"sections"`
	TrackingParams string                         `json:"trackingParams"`
	Style          *string                        `json:"style,omitempty"`
}

type MultiPageMenuRendererSection struct {
	MultiPageMenuSectionRenderer MultiPageMenuSectionRenderer `json:"multiPageMenuSectionRenderer"`
}

type MultiPageMenuSectionRenderer struct {
	Items          []MultiPageMenuSectionRendererItem `json:"items"`
	TrackingParams string                             `json:"trackingParams"`
}

type MultiPageMenuSectionRendererItem struct {
	CompactLinkRenderer CompactLinkRenderer `json:"compactLinkRenderer"`
}

type CompactLinkRenderer struct {
	Icon               IconImage                             `json:"icon"`
	Title              Text                                  `json:"title"`
	NavigationEndpoint CompactLinkRendererNavigationEndpoint `json:"navigationEndpoint"`
	TrackingParams     string                                `json:"trackingParams"`
	Style              *string                               `json:"style,omitempty"`
}

type CompactLinkRendererNavigationEndpoint struct {
	ClickTrackingParams      string                `json:"clickTrackingParams"`
	CommandMetadata          PurpleCommandMetadata `json:"commandMetadata"`
	UploadEndpoint           *UploadEndpoint       `json:"uploadEndpoint,omitempty"`
	SignalNavigationEndpoint *Signal               `json:"signalNavigationEndpoint,omitempty"`
	URLEndpoint              *URLEndpoint          `json:"urlEndpoint,omitempty"`
}

type PurpleCommandMetadata struct {
	WebCommandMetadata StickyWebCommandMetadata `json:"webCommandMetadata"`
}

type StickyWebCommandMetadata struct {
	URL    string `json:"url"`
	RootVe int64  `json:"rootVe"`
}

type URLEndpoint struct {
	URL    string `json:"url"`
	Target string `json:"target"`
}

type UploadEndpoint struct {
	Hack bool `json:"hack"`
}

type URL string

const (
	ServiceAjax URL = "/service_ajax"
)

type WebPageType string

const (
	WebPageTypeBrowse WebPageType = "WEB_PAGE_TYPE_BROWSE"
	WebPageTypeSearch WebPageType = "WEB_PAGE_TYPE_SEARCH"
	WebPageTypeWatch  WebPageType = "WEB_PAGE_TYPE_WATCH"
)

type ID string

const (
	UCVk4G4NEtBS4TLOyHqustDA ID = "UCVk4G4nEtBS4tLOyHqustDA"
)

type CanonicalBaseURL string

const (
	UserStatacorp CanonicalBaseURL = "/user/statacorp"
)
