package main

/*
#include <stdlib.h>
*/
import "C"
import (
    "bufio"
    "encoding/json"
    "strings"
    "sync"
    "time"
)

type Contact struct {
    Organization *string `json:"organization,omitempty"`
    Email        *string `json:"email,omitempty"`
    Name         *string `json:"name,omitempty"`
    Telephone    *string `json:"telephone,omitempty"`
}

type Tech Contact
type Registrant Contact
type Admin Contact

type Abuse struct {
    Email     *string `json:"email,omitempty"`
    Telephone *string `json:"telephone,omitempty"`
}

type WhoisRecord struct {
    RawText string `json:"raw_text"`

    Registrant    Registrant `json:"registrant"`
    Admin         Admin      `json:"admin"`
    Tech          Tech       `json:"tech"`
    Abuse         Abuse      `json:"abuse"`
    Statuses      []string   `json:"statuses"`
    NameServers   []string   `json:"name_servers"`
    Domain        *string    `json:"domain,omitempty"`
    Registrar     *string    `json:"registrar,omitempty"`
    ExpiresAt     *int64     `json:"expires_at,omitempty"`
    RegisteredAt  *int64     `json:"registered_at,omitempty"`
    UpdatedAt     *int64     `json:"updated_at,omitempty"`
    IsRateLimited bool       `json:"is_rate_limited"`
}

//export ParseWhois
func ParseWhois(rawText *C.char) *C.char {
    normalizedRawText := normalizeRawText(C.GoString(rawText))
    wr := WhoisRecord{
       RawText: normalizedRawText,
    }

    var wg sync.WaitGroup
    var mu sync.Mutex
    wg.Add(12)

    go func() {
       defer wg.Done()
       abuse := findAbuse(normalizedRawText)
       mu.Lock()
       wr.Abuse = abuse
       mu.Unlock()
    }()
    go func() {
       defer wg.Done()
       admin := findAdmin(normalizedRawText)
       mu.Lock()
       wr.Admin = admin
       mu.Unlock()
    }()
    go func() {
       defer wg.Done()
       domain := findDomain(normalizedRawText)
       mu.Lock()
       wr.Domain = domain
       mu.Unlock()
    }()
    go func() {
       defer wg.Done()
       expiresAt := findExpiresAt(normalizedRawText)
       mu.Lock()
       wr.ExpiresAt = expiresAt
       mu.Unlock()
    }()
    go func() {
       defer wg.Done()
       nameServers := findNameServers(normalizedRawText)
       mu.Lock()
       wr.NameServers = nameServers
       mu.Unlock()
    }()
    go func() {
       defer wg.Done()
       registeredAt := findRegisteredAt(normalizedRawText)
       mu.Lock()
       wr.RegisteredAt = registeredAt
       mu.Unlock()
    }()
    go func() {
       defer wg.Done()
       updatedAt := findUpdatedAt(normalizedRawText)
       mu.Lock()
       wr.UpdatedAt = updatedAt
       mu.Unlock()
    }()
    go func() {
       defer wg.Done()
       registrant := findRegistrant(normalizedRawText)
       mu.Lock()
       wr.Registrant = registrant
       mu.Unlock()
    }()
    go func() {
       defer wg.Done()
       registrar := findRegistrar(normalizedRawText)
       mu.Lock()
       wr.Registrar = registrar
       mu.Unlock()
    }()
    go func() {
       defer wg.Done()
       statuses := findStatuses(normalizedRawText)
       mu.Lock()
       wr.Statuses = statuses
       mu.Unlock()
    }()
    go func() {
       defer wg.Done()
       tech := findTech(normalizedRawText)
       mu.Lock()
       wr.Tech = tech
       mu.Unlock()
    }()
    go func() {
       defer wg.Done()
       isRateLimited := findIsRateLimited(normalizedRawText)
       mu.Lock()
       wr.IsRateLimited = isRateLimited
       mu.Unlock()
    }()

    wg.Wait()

    jsonData, err := json.MarshalIndent(wr, "", "    ")
    if err != nil {
       return C.CString("")
    }
    return C.CString(string(jsonData))
}

func main() {}

func normalizeRawText(rawText string) string {
    scanner := bufio.NewScanner(strings.NewReader(rawText))
    var lines []string
    for scanner.Scan() {
       line := strings.TrimSpace(scanner.Text())
       lines = append(lines, line)
    }
    reversedLines := reverse(lines)
    sharpIndex := -1
    for i, line := range reversedLines {
       if strings.HasPrefix(line, "#") {
          sharpIndex = i + 1
          break
       }
    }
    if sharpIndex == -1 {
       return strings.Join(lines, "\n")
    }
    lines = reversedLines[:sharpIndex]
    lines = reverse(lines)
    return strings.Join(lines, "\n")
}

func reverse(s []string) []string {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
       s[i], s[j] = s[j], s[i]
    }
    return s
}

func findByKeywords(normalizedRawText string, keywords []string) *string {
    lines := strings.Split(normalizedRawText, "\n")
    for _, line := range lines {
       for _, keyword := range keywords {
          if strings.HasPrefix(strings.ToLower(line), strings.ToLower(keyword)) {
             parts := strings.SplitN(line, ":", 2)
             if len(parts) > 1 {
                value := strings.TrimSpace(parts[1])
                return &value
             }
          }
       }
    }
    return nil
}

func findAllByKeywords(normalizedRawText string, keywords []string) []string {
    lines := strings.Split(normalizedRawText, "\n")
    var results []string
    for _, line := range lines {
       for _, keyword := range keywords {
          if strings.HasPrefix(strings.ToLower(line), strings.ToLower(keyword)) {
             parts := strings.SplitN(line, ":", 2)
             if len(parts) > 1 {
                value := strings.TrimSpace(parts[1])
                results = append(results, value)
             }
          }
       }
    }
    return results
}

func findDateByKeywords(normalizedRawText string, keywords []string, dateFormats []string) *int64 {
   lines := strings.Split(normalizedRawText, "\n")
   for _, line := range lines {
       for _, keyword := range keywords {
           if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
               parts := strings.SplitN(line, ":", 2)
               if len(parts) > 1 {
                   dateStr := strings.TrimSpace(parts[1])
                   var parsedTime time.Time
                   var err error
                   for _, format := range dateFormats {
                       parsedTime, err = time.Parse(format, dateStr)
                       if err == nil {
                           unixTimestamp := parsedTime.Unix()
                           return &unixTimestamp
                       }
                   }
               }
           }
       }
   }

   return nil
}

func findAbuseEmail(normalizedRawText string) *string {
    keywords := []string{
       "Registrar Abuse Contact Email",
       "AC E-Mail",
    }
    return findByKeywords(normalizedRawText, keywords)
}

func findAbuseTelephone(normalizedRawText string) *string {
    keywords := []string{
       "Registrar Abuse Contact Phone",
       "AC Phone Number",
    }
    return findByKeywords(normalizedRawText, keywords)
}

func findAbuse(normalizedRawText string) Abuse {
    return Abuse{
       Email:     findAbuseEmail(normalizedRawText),
       Telephone: findAbuseTelephone(normalizedRawText),
    }
}

func findAdminName(normalizedRawText string) *string {
    keywords := []string{"Admin Name"}
    return findByKeywords(normalizedRawText, keywords)
}

func findAdminEmail(normalizedRawText string) *string {
    keywords := []string{"Admin Email"}
    return findByKeywords(normalizedRawText, keywords)
}

func findAdminTelephone(normalizedRawText string) *string {
    keywords := []string{"Admin Phone", "Admin Phone:"}
    return findByKeywords(normalizedRawText, keywords)
}

func findAdminOrganization(normalizedRawText string) *string {
    keywords := []string{"Admin Organization"}
    return findByKeywords(normalizedRawText, keywords)
}

func findAdmin(normalizedRawText string) Admin {
    return Admin{
       Name:         findAdminName(normalizedRawText),
       Email:        findAdminEmail(normalizedRawText),
       Telephone:    findAdminTelephone(normalizedRawText),
       Organization: findAdminOrganization(normalizedRawText),
    }
}

func findDomain(normalizedRawText string) *string {
    keywords := []string{"Domain Name", "domain:"}
    return findByKeywords(normalizedRawText, keywords)
}

func findNameServers(normalizedRawText string) []string {
   keywords := []string{"Name server", "Nameserver", "nameservers", "Nserver", "Host Name"}
   values := findAllByKeywords(normalizedRawText, keywords)

   // Use a map to track unique name servers
   nameServerMap := make(map[string]struct{})
   for _, value := range values {
       nameServer := strings.ToLower(strings.Split(value, " ")[0])
       nameServerMap[nameServer] = struct{}{}
   }

   // Convert the map values to a slice
   var nameServers []string
   for nameServer := range nameServerMap {
       nameServers = append(nameServers, nameServer)
   }

   return nameServers
}

func findUpdatedAt(normalizedRawText string) *int64 {
   keywords := []string{
       "Updated Date", "Update Date", "updated", "changed", "modified",
       "Last Updated On", "Last Updated Date", "domain_datelastmodified",
       "Last Update", "Modified Date", "last-update",
   }
   dateFormats := []string{
       time.RFC3339, "2006-01-02T15:04:05Z", "2006-01-02", 
       "January 2, 2006", "02-Jan-2006",
   }
   return findDateByKeywords(normalizedRawText, keywords, dateFormats)
}

func findExpiresAt(normalizedRawText string) *int64 {
   keywords := []string{
       "Expiry Date", "Expiration Date", "Expire Date", "expire", "expires",
       "Expires On", "Expiration Time", "Renewal Date", "Record expires on",
       "paid-till", "expire-date", "domain_datebilleduntil", "Valid Until", 
       "validity",
   }
   dateFormats := []string{
       time.RFC3339, "2006-01-02T15:04:05Z", "2006-01-02", 
       "January 2, 2006", "02-Jan-2006",
   }
   return findDateByKeywords(normalizedRawText, keywords, dateFormats)
}

func findRegisteredAt(normalizedRawText string) *int64 {
   keywords := []string{
       "Creation Date", "Creation Date:", "registered", "created", "activated",
       "Registration Time", "Registered Date", "Registration Date", 
       "Record created on", "Created On", "registered on", "Created Date",
   }
   dateFormats := []string{
       time.RFC3339, "2006-01-02T15:04:05Z", "2006-01-02", 
       "January 2, 2006", "02-Jan-2006",
   }
   return findDateByKeywords(normalizedRawText, keywords, dateFormats)
}

func findRegistrantName(normalizedRawText string) *string {
    keywords := []string{
       "Registrant Name",
       "Registrant",
       "Registrant Contact Name",
       "Person",
       "registrant_contact_name",
       "Domain Holder",
       "personname",
       "responsible",
    }
    return findByKeywords(normalizedRawText, keywords)
}

func findRegistrantEmail(normalizedRawText string) *string {
    keywords := []string{
       "Registrant Email",
       "Registrant Contact Email",
    }
    return findByKeywords(normalizedRawText, keywords)
}

func findRegistrantTelephone(normalizedRawText string) *string {
    keywords := []string{
       "Registrant Phone","Registrant Phone:",
    }
    return findByKeywords(normalizedRawText, keywords)
}

func findRegistrantOrganization(normalizedRawText string) *string {
    keywords := []string{
       "Registrant Organization",
       "org",
       "org-name",
       "Registrant Contact Organisation",
       "Domain Holder Organization",
    }
    return findByKeywords(normalizedRawText, keywords)
}

func findRegistrant(normalizedRawText string) Registrant {
    return Registrant{
       Name:         findRegistrantName(normalizedRawText),
       Email:        findRegistrantEmail(normalizedRawText),
       Telephone:    findRegistrantTelephone(normalizedRawText),
       Organization: findRegistrantOrganization(normalizedRawText),
    }
}

func findRegistrar(normalizedRawText string) *string {
    keywords := []string{
       "Registrar:",
       "Registrar Name",
       "Sponsoring Registrar",
       "registrar-name",
       "Registration Service Provider",
       "Domain Support",
       "Sponsoring Registrar Organization",
       "Account Name",
    }

    return findByKeywords(normalizedRawText, keywords)
}

func findStatuses(normalizedRawText string) []string {
   keywords := []string{"Domain Status", "domaintype"}
   values := findAllByKeywords(normalizedRawText, keywords)

   // Use a map to track unique statuses
   statusMap := make(map[string]struct{})
   for _, value := range values {
       if value == "No Object Found" {
           statusMap[value] = struct{}{}
       } else {
           parts := strings.Fields(value)
           if len(parts) > 0 {
               statusMap[parts[0]] = struct{}{}
           }
       }
   }

   // Convert the map values to a slice
   var statuses []string
   for status := range statusMap {
       statuses = append(statuses, status)
   }

   return statuses
}

func findTechName(normalizedRawText string) *string {
    keywords := []string{"Tech Name", "Tech Contact Name", "Tech Contact"}
    return findByKeywords(normalizedRawText, keywords)
}

func findTechEmail(normalizedRawText string) *string {
    keywords := []string{"Tech Email", "Tech Contact Email"}
    return findByKeywords(normalizedRawText, keywords)
}

func findTechTelephone(normalizedRawText string) *string {
    keywords := []string{"Tech Phone"}
    return findByKeywords(normalizedRawText, keywords)
}

func findTechOrganization(normalizedRawText string) *string {
    keywords := []string{"Tech Organization", "Tech Contact Organisation"}
    return findByKeywords(normalizedRawText, keywords)
}

func findTech(normalizedRawText string) Tech {
    return Tech{
       Name:         findTechName(normalizedRawText),
       Email:        findTechEmail(normalizedRawText),
       Telephone:    findTechTelephone(normalizedRawText),
       Organization: findTechOrganization(normalizedRawText),
    }
}

func findIsRateLimited(rawText string) bool {
    rateLimitStrings := []string{
       "WHOIS LIMIT EXCEEDED - SEE WWW.PIR.ORG/WHOIS FOR DETAILS",
       "Your access is too fast,please try again later.",
       "Your connection limit exceeded.",
       "Number of allowed queries exceeded.",
       "WHOIS LIMIT EXCEEDED",
       "Requests of this client are not permitted.",
       "Too many connection attempts. Please try again in a few seconds.",
       "We are unable to process your request at this time.",
       "HTTP/1.1 400 Bad Request",
       "Closing connections because of Timeout",
       "Access to whois service at whois.isoc.org.il was **DENIED**",
       "IP Address Has Reached Rate Limit",
    }

    rawText = strings.TrimSpace(rawText)

    for _, limitString := range rateLimitStrings {
       if strings.Contains(rawText, limitString) {
          return true
       }
    }

    return false
}