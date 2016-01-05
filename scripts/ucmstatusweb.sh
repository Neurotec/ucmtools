#!/bin/bash
#crea una pagina html con el estado de los ucm indicados
UCMSTATUS=../tools/ucmstatus


#cabezera
cat <<EOF
<!DOCTYPE html>
<html>
<head>
<title>UCM STATUS</title>
</head>
<body>
<table border="1">
<thead>
<tr>
 <th>UCM</th>
 <th>PROFILE</th>
 <th>TRUNK</th>
 <th>USER</th>
 <th>TYPE</th>
 <th>STATUS</th>
</tr>
</thead>
<tbody>
EOF

while read line
do
    hostname=$(echo $line | cut -d' ' -f1)
    host=$(echo $line | cut -d' ' -f2)
    username=$(echo $line | cut -d' ' -f3)
    password=$(echo $line | cut -d' ' -f4)
    IFS=$'\n'
    for ucmstatus in $(env UCMPASSWORD=$password $UCMSTATUS $host $username |  tr -t '\t' '|' | tr -s '|')
    do
	status=$(echo -n $ucmstatus | cut -d'|' -f5)
	if [[ $status == "Registered" || $status == "Idle" ]]; then
	    echo -n "<tr style='color: green'>"
	else
	    echo -n "<tr style='color: red'>"
	fi

	echo -n "<td>"
	echo -n "<a target='__blank' href='https://$host'>$hostname</a>"
	echo -n "</td>"
	echo -n "<td>"
	echo -n $ucmstatus | cut -d'|' -f1
	echo -n "</td>"
	echo -n "<td>"
	echo -n $ucmstatus | cut -d'|' -f2
	echo -n "</td>"
	echo -n "<td>"
	echo -n $ucmstatus | cut -d'|' -f3
	echo -n "</td>"
	echo -n "<td>"
	echo -n $ucmstatus | cut -d'|' -f4
	echo -n "</td>"
	echo -n "<td>"
	echo -n $status
	echo -n "</td>"
	echo -n "</tr>"
    done
done < ./ucmstatusweb.conf

cat <<EOF
</tbody>
</table>
<bold>ULTIMA ACTUALIZACION: $(date)</bold>
</body>
EOF
