package templates

var InvitationTemplate = `
<!DOCTYPE html>
<html>
<head>

</head>
<body style="background-color: #1C444E; color: white;">
<div style="display: flex; justify-content: center; flex-direction: column; max-width: 30rem; margin: auto; font-family: Georgia;">
	<img src="https://cdn.logo.com/hotlink-ok/logo-social-sq.png" width="120" alt="logo">
    <h1 style="font-size: 42px;">Invitation to the Party</h1>
    <h3>{{.GuestName}} is invited to join the annual Christmas Party</h3>
    <h3>Reserved Table Number {{.TableNumber}}</h3>
    <h3>with {{.AccompanyingGuests}} accompanying guests</h3>
</div>
</body>
</html>
`
