var dateRelativityApplicator = function($selector) {
  $selector.find("time[datetime]").each(function(index, el) {
    var $el = $(el),
    m = moment($el.attr("datetime"));

    $el.attr("title", $el.attr("datetime"));
    if ($el.data("role") == "timeago") {
      $el.text(m.fromNow());
    } else {
      $el.text(m.calendar());
    }
  });
}

$(function() {
  dateRelativityApplicator($("body"));

  $resultsContainer = $("#results-container")
  if ($resultsContainer.length) {
    var source = $("#results-template").html();
    var template = Handlebars.compile(source);
    var url = $resultsContainer.data("url");

    $resultsContainer.text("Loading…")
    $.getJSON(url, function(data, status, xhr) {
      $.map(data, function(item, index) {
        if (item.Success) {
          item.Icon = "glyphicon-ok";
          item.TextStatus = "OK"
          item.CSSclass = "default"
        } else if (item.status < 100) {
          item.Icon = "glyphicon-warning-sign";
          item.TextStatus = "Warning"
          item.CSSclass = "warning"
        } else {
          item.Icon = "glyphicon-fire";
          item.TextStatus = "Fail!"
          item.CSSclass = "danger"
        }
        return item
      });

      var html = template({"Results": data});
      $resultsContainer.html(html);
      dateRelativityApplicator($resultsContainer);
    });
  }
});
