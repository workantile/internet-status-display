package main

import (
	"encoding/json"
	"html/template"
	"log"
	"math"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Status struct {
	mutex sync.RWMutex `json:"-"`

	FastcomSpeedKbps int        `json:"fastcomKbps"`
	GooglePing       *PingResult `json:"google"`
	CloudflarePing   *PingResult `json:"cloudflare"`

	RouterPing       *PingResult `json:"router"`
	SwitchPing       *PingResult `json:"switch"`
	CloudKeyPing     *PingResult `json:"cloudKey"`
	DownstairsAPPing *PingResult `json:"downstairsAP"`
	LoftAPPing       *PingResult `json:"loftAP"`
	PhoneRoomsAPPing *PingResult `json:"phoneRoomsAP"`
}

func (s *Status) Update(updater func(s *Status)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	updater(s)
}

var CurrentStatus = Status{}
var ProbeMutex sync.Mutex

func statusHandler(w http.ResponseWriter, r *http.Request) {
	CurrentStatus.mutex.RLock()
	defer CurrentStatus.mutex.RUnlock()

	j, err := json.Marshal(CurrentStatus)
	if err != nil {
		http.Error(w, errors.Wrap(err, "marshal current status to JSON").Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(j)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, errors.Wrap(err, "read template").Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, struct{}{}); err != nil {
		http.Error(w, errors.Wrap(err, "render template").Error(), http.StatusInternalServerError)
	}
}

const pingInterval = 10 * time.Second
const googleAddr = "8.8.8.8"
const cloudflareAddr = "1.1.1.1"
const routerAddr = "10.10.10.1"
const switchAddr = "10.10.10.255"  // TODO(cdzombak):
const cloudKeyAddr = "10.10.10.14"
const downstairsAPAddr = "10.10.10.255"  // TODO(cdzombak):
const loftAPAddr = "10.10.10.255"  // TODO(cdzombak):
const phoneRoomsAPAddr = "10.10.10.255"  // TODO(cdzombak):

func main() {
	startPing(pingInterval, googleAddr, func(stats *PingResult, err error) {
		if err != nil {
			log.Println("[warning] google ping:", err)
			CurrentStatus.Update(func(s *Status) {
				s.GooglePing = nil
			})
			return
		}
		log.Printf("[debug] google ping avg: %d ms\n", stats.AvgRtt)
		log.Printf("[debug] google pkt loss: %d%%\n", stats.PacketLossPct)
		CurrentStatus.Update(func(s *Status) {
			s.GooglePing = stats
		})
	})

	startPing(pingInterval, cloudflareAddr, func(stats *PingResult, err error) {
		if err != nil {
			log.Println("[warning] cloudflare ping:", err)
			CurrentStatus.Update(func(s *Status) {
				s.CloudflarePing = nil
			})
			return
		}
		log.Printf("[debug] cloudflare ping avg: %d ms\n", stats.AvgRtt)
		log.Printf("[debug] cloudflare pkt loss: %d%%\n", stats.PacketLossPct)
		CurrentStatus.Update(func(s *Status) {
			s.CloudflarePing = stats
		})
	})

	startPing(pingInterval, routerAddr, func(stats *PingResult, err error) {
		if err != nil {
			log.Println("[warning] router ping:", err)
			CurrentStatus.Update(func(s *Status) {
				s.RouterPing = nil
			})
			return
		}
		log.Printf("[debug] router ping avg: %d ms\n", stats.AvgRtt)
		log.Printf("[debug] router pkt loss: %d%%\n", stats.PacketLossPct)
		CurrentStatus.Update(func(s *Status) {
			s.RouterPing = stats
		})
	})

	startPing(pingInterval, cloudKeyAddr, func(stats *PingResult, err error) {
		if err != nil {
			log.Println("[warning] controller ping:", err)
			CurrentStatus.Update(func(s *Status) {
				s.CloudKeyPing = nil
			})
			return
		}
		log.Printf("[debug] controller ping avg: %d ms\n", stats.AvgRtt)
		log.Printf("[debug] controller pkt loss: %d%%\n", stats.PacketLossPct)
		CurrentStatus.Update(func(s *Status) {
			s.CloudKeyPing = stats
		})
	})

	startPing(pingInterval, switchAddr, func(stats *PingResult, err error) {
		if err != nil {
			log.Println("[warning] switch ping:", err)
			CurrentStatus.Update(func(s *Status) {
				s.SwitchPing = nil
			})
			return
		}
		log.Printf("[debug] switch ping avg: %d ms\n", stats.AvgRtt)
		log.Printf("[debug] switch pkt loss: %d%%\n", stats.PacketLossPct)
		CurrentStatus.Update(func(s *Status) {
			s.SwitchPing = stats
		})
	})

	startPing(pingInterval, downstairsAPAddr, func(stats *PingResult, err error) {
		if err != nil {
			log.Println("[warning] downstairs AP ping:", err)
			CurrentStatus.Update(func(s *Status) {
				s.DownstairsAPPing = nil
			})
			return
		}
		log.Printf("[debug] downstairs AP ping avg: %d ms\n", stats.AvgRtt)
		log.Printf("[debug] downstairs AP pkt loss: %d%%\n", stats.PacketLossPct)
		CurrentStatus.Update(func(s *Status) {
			s.DownstairsAPPing = stats
		})
	})

	startPing(pingInterval, loftAPAddr, func(stats *PingResult, err error) {
		if err != nil {
			log.Println("[warning] loft AP ping:", err)
			CurrentStatus.Update(func(s *Status) {
				s.LoftAPPing = nil
			})
			return
		}
		log.Printf("[debug] loft AP ping avg: %d ms\n", stats.AvgRtt)
		log.Printf("[debug] loft AP pkt loss: %d%%\n", stats.PacketLossPct)
		CurrentStatus.Update(func(s *Status) {
			s.LoftAPPing = stats
		})
	})

	startPing(pingInterval, phoneRoomsAPAddr, func(stats *PingResult, err error) {
		if err != nil {
			log.Println("[warning] phone rooms AP ping:", err)
			CurrentStatus.Update(func(s *Status) {
				s.PhoneRoomsAPPing = nil
			})
			return
		}
		log.Printf("[debug] phone rooms AP ping avg: %d ms\n", stats.AvgRtt)
		log.Printf("[debug] phone rooms AP pkt loss: %d%%\n", stats.PacketLossPct)
		CurrentStatus.Update(func(s *Status) {
			s.PhoneRoomsAPPing = stats
		})
	})

	startFastcomSpeedTest(4*time.Minute, func(kbps float64, err error) {
		if err != nil {
			log.Println("[warning] speedtest:", err)
			CurrentStatus.Update(func(s *Status) {
				s.FastcomSpeedKbps = 0
			})
			return
		}
		log.Printf("[debug] speedtest result: %.2f Kbps\n", kbps)
		CurrentStatus.Update(func(s *Status) {
			s.FastcomSpeedKbps = int(math.Round(kbps))
		})
	})

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/status", statusHandler)
	log.Fatal(http.ListenAndServe(":8001", nil))
}