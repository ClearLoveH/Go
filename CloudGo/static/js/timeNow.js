$(document).ready(function() {
    $.ajax({
        url: "/today"
    }).then(function(data) {
       $('.time').append(data);
    });
});
