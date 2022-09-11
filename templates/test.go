package templates

{{define "content"}}
<div class="col-md-12">
{{$res := index .Data "reservations"}}

<table class="table table-striped table-hover" id="all-res">
<thead>
<tr>
<th>ID</th>
<th>Last Name</th>
<th>Room</th>
<th>Arrival</th>
<th>Departure</th>
</tr>
</thead>
<tbody>
{{range $res}}
<td>{{.ID}}</td>
<td>
<a href="/bookings/admin/reservations/all/{{.ID}}/show">
{{.LastName}}
</a>
</td>
<td>{{.Room.RoomName}}</td>
<td>{{humanDate .StartDate}}</td>
<td>{{humanDate .EndDate}}</td>
</tbody>


{{end}}
</table>

</div>
{{end}}