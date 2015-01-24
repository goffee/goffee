$(function() {
  $("time[datetime]").each(function(index, el) {
    var $el = $(el),
        m = moment($el.attr("datetime"));

    $el.attr("title", $el.attr("datetime"));
    if ($el.data("role") == "timeago") {
      $el.text(m.fromNow());
    } else {
      $el.text(m.calendar());
    }
  });
});
