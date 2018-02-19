require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap-sass/assets/javascripts/bootstrap.js");
$(() => {

});

$(document).on('change', '#fish-label', function() {
  console.log($('#fish-label'));
  $.ajax({
    type: "POST",
    url: '/label',
    success: function (data) {
      console.log(data);
      location.reload();
    },
    data: {'fish_pic_id': 123,
           'fish_label': $('#fish-label')[0].value}
  });
});

$("#fish-pic-form").change(function() {
  $('#fish-pic-form').submit();
});
