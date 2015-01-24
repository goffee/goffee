$(function() {
  $("time[datetime]").each(function(index, el) {
    var $el;
    $el = $(el);
    $el.attr("title", $el.attr("datetime"))
    return $el.text(moment($el.attr("datetime")).fromNow());
  });
});
