$(document).ready(function() {
    $.ajax({
        url: "/timeNow"
    }).then(function(data) {
       $('.time').append(data);
    });
});
