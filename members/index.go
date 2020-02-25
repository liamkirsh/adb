package members

import (
	"database/sql"
	"fmt"
	"html/template"
	"sort"
)

// TODO(mdempsky): Use adb_users instead?
var adminEmails = map[string]bool{
	"matthew@dempsky.org": true,
}

func (s *server) index() {
	email, err := s.googleEmail()
	if err != nil {
		s.redirect(absURL("/login"))
		return
	}

	if adminEmails[email] {
		if q := s.r.URL.Query()["email"]; len(q) >= 1 && q[0] != "" {
			email = q[0]
		}
	}

	// MySQL doesn't have a proper boolean data type, and it's
	// json_object seems to have some arbitrary heuristics for
	// deciding when to encode a boolean expression as 0/1 vs
	// true/false.
	var data struct {
		Name          string
		Email         string
		Phone         string
		Location      string
		Facebook      string
		Birthday      string
		ActivistLevel string

		Organizer     bool
		ChapterMember bool
		FebVoter      bool
		MarVoter      bool

		WorkingGroups []string

		Total      int
		Attendance []struct {
			Month        int // YYYYMM
			MPI          int // boolean
			Community    int // boolean
			DirectAction int // boolean
			Events       []struct {
				Date         string // "YYYY-MM-DD"
				Name         string
				Community    int // boolean
				DirectAction int //  boolean
			}
		}
	}

	// This query would be more natural if attendance could be
	// computed using a subquery like working groups, but because
	// of the two-level aggregation, we'd actually need a
	// sub-subquery; and subqueries can only access variables from
	// the immediately outer context.
	const q = `
select json_object(
  'Name', x.name,
  'Email', x.email,
  'Phone', x.phone,
  'Location', x.location,
  'Facebook', x.facebook,
  'Birthday', x.dob,
  'ActivistLevel', x.activist_level,

  'Organizer', x.activist_level in ('Organizer', 'Senior Organizer'),
  'ChapterMember', x.activist_level in ('Chapter Member', 'Organizer', 'Senior Organizer'),

  'FebVoter', sum(x.mpi and x.month >= 201911 and x.month < 202002) >= 2,
  'MarVoter', sum(x.mpi and x.month >= 201912 and x.month < 202003) >= 2,

  'WorkingGroups', (
    select json_arrayagg(w.name)
    from working_groups w
    join working_group_members m on (w.id = m.working_group_id)
    where m.activist_id = x.id
  ),

  'Total', sum(x.subtotal),
  'Attendance', if(sum(x.subtotal) = 0, null,
    json_arrayagg(json_object(
      'Month', x.month,
      'MPI', x.mpi,
      'Community', x.community,
      'DirectAction', x.direct_action,
      'Events', x.events
    )))
)
from (
  select a.id, a.name, a.email, a.phone, a.location, a.facebook, a.activist_level, a.dob, a.date_organizer,
    e.month, count(e.id) as subtotal,
    max(e.community) as community, max(e.direct_action) as direct_action,
    (max(e.direct_action) and (max(e.community) or (e.month >= 202001 and e.month < 202003))) as mpi,
    json_arrayagg(json_object(
      'Date', e.date,
      'Name', e.name,
      'Community', e.community,
      'DirectAction', e.direct_action
    )) as events
  from activists a
  left join event_attendance ea on (a.id = ea.activist_id)
  left join (
          select id, date,
                 concat(name, if(event_type = 'Connection', ' (Connection)', '')) as name,
                 extract(year_month from date) as month,
                 event_type in ('Circle', 'Community', 'Training') as community,
                 event_type in ('Action', 'Campaign Action', 'Frontline Surveillance', 'Outreach', 'Sanctuary') as direct_action
          from events
        ) e on (e.id = ea.event_id)
  where a.email = ?
    and not a.hidden
  group by a.id, e.month
) x
group by x.id
`

	if err := s.queryJSON(&data, q, email); err != nil {
		if err == sql.ErrNoRows {
			s.render(absentTmpl, email)
		} else {
			s.error(err)
		}
		return
	}

	// Manually sort in descending order by date, as MySQL doesn't
	// allow control of json_arrayagg()'s aggregation order.
	sort.Slice(data.Attendance, func(i, j int) bool { return data.Attendance[i].Month > data.Attendance[j].Month })
	for k := range data.Attendance {
		events := data.Attendance[k].Events
		sort.Slice(events, func(i, j int) bool { return events[i].Date > events[j].Date })
	}

	s.render(indexTmpl, &data)
}

var absentTmpl = template.Must(template.New("absent").Parse(`
<!doctype html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link href="https://fonts.googleapis.com/css?family=Source+Sans+Pro&display=swap" rel="stylesheet">

<style>
body {
  font-family: 'Source Sans Pro', sans-serif;
}

.wrap {
  max-width: 40em;
  margin-left: auto;
  margin-right: auto;
}
</style>
</head>

<body>
<div class="wrap">
<p>Sorry, we don't have any records associated with the email address "{{.}}".</p>
<p>You can try <a href="login?force">logging in with another account</a>
or email <a href="mailto:tech@dxe.io">tech@dxe.io</a> for help.</p>
</div>
</body>
`))

var indexTmpl = template.Must(template.New("index").Funcs(template.FuncMap{
	"monthfmt": func(n int) string { return fmt.Sprintf("%d-%02d", n/100, n%100) },
}).Parse(`
<!doctype html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link href="https://fonts.googleapis.com/css?family=Source+Sans+Pro&display=swap" rel="stylesheet">

<style>
body {
  font-family: 'Source Sans Pro', sans-serif;
}

.wrap {
  max-width: 40em;
  margin-left: auto;
  margin-right: auto;
}

h1, h2 {
  margin-top: 2em;
}

table {
  border-spacing: 0;
  font-variant-numeric: tabular-nums;
}

table.attendance {
  padding-top: 2em;
}

tr.month {
  background-color: #ddd;
  font-weight: bold;
}

tr.month td {
  text-align: center;
}

tr.mpi {
  background-color: #beb;
}

td {
  padding: 0.375em;
}

table.attendance td:nth-child(3) {
  white-space: nowrap;
}

table.profile td:nth-child(1), table.election, td:nth-child(1) {
  font-weight: bold;
}

.green { background-color: #beb; }
.gray { background-color: #ddd; }
</style>
</head>

<body>
<div class="wrap">

<h1>Activist Record</h1>

<p>Hello, <b>{{.Name}}</b>! (Not you? <a href="login?force">Click here</a> to login as someone else.)</p>

<p>Questions or feedback about this page? Email <a href="mailto:tech@dxe.io">tech@dxe.io</a>.</p>

<h2>Profile</h2>

<table class="profile">
<tr><td>Name:</td><td>{{.Name}}</td></tr>
<tr><td>Email:</td><td>{{.Email}}</td></tr>
<tr><td>Phone:</td><td>{{.Phone}}</td></tr>
<tr><td>Location:</td><td>{{.Location}}</td></tr>
<tr><td>Facebook Profile:</td><td><a href="{{.Facebook}}">{{.Facebook}}</a></td></tr>
<tr><td>Birthday:</td><td>{{.Birthday}}</td></tr>
<tr><td><a href="https://docs.google.com/document/d/1QnJXz8YuQeBL0cz4iK60mOvQfDN1vd7SBwvVhRFDHNc/preview">Activist Level</a>:</td><td>{{.ActivistLevel}}</td></tr>
</table>

<h2>Voter Eligibility</h2>

<table class="elections">
<tr>
  <th>Month</th>
  <th><a href="https://docs.dxesf.org/#32-chapter-members">Chapter Member Votes</a></th>
  <th><a href="https://docs.dxesf.org/#33-organizers">Organizer Votes</a></th>
</tr>
<tr>
  <td>February 2020</td>
  <td>{{if (.ChapterMember) (.FebVoter)}}Yes{{else}}No{{end}}</td>
  <td>{{if (.Organizer) (.FebVoter)}}Yes{{else}}No{{end}}</td>
</tr>
<tr>
  <td>March 2020</td>
  <td>{{if (.ChapterMember) (.MarVoter)}}Yes{{else}}No{{end}}</td>
  <td>{{if (.Organizer) (.MarVoter)}}Yes{{else}}No{{end}}</td>
</tr>
</table>

<h2>Working Groups</h2>

{{if .WorkingGroups}}
<ul>
{{range .WorkingGroups}}
<li>{{.}}</li>
{{end}}
</ul>
{{else}}
<p>None.</p>
{{end}}

<h2>Event Attendance</h2>

<p>Below are <b>{{.Total}}</b> events you've attended with DxE SF.</p>

<p>🏙️ indicates a "community" event;<br>
📣 indicates a "direct action" event.</p>

<p>A <b class="green">green</b> bar indicates you met MPI requirements that month;<br>
a <b class="gray">gray</b> bar indicates you did not.</p>

<table class="attendance">
{{range .Attendance}}
<tr class="month {{if .MPI}}mpi{{end}}">
  <td>{{if .Community}}🏙️{{else if (ge .Month 202001) (le .Month 202002)}}🆓{{end}}</td>
  <td>{{if .DirectAction}}📣{{end}}</td>
  <td colspan=2>{{monthfmt .Month}}</td>
</tr>
{{range .Events}}
<tr>
  <td>{{if .Community}}🏙️{{end}}</td>
  <td>{{if .DirectAction}}📣{{end}}</td>
  <td>{{.Date}}</td>
  <td>{{.Name}}</td>
</tr>
{{end}}
{{end}}
</table>

</div>
</body>
</html>
`))
