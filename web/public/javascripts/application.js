//
// Regular Expression for URL validation
//
// Author: Diego Perini
// Updated: 2010/12/05
// License: MIT
//
// Copyright (c) 2010-2013 Diego Perini (http://www.iport.it)
// https://gist.github.com/dperini/729294
var re_weburl = new RegExp(
  "^" +
  // protocol identifier
  "(?:(?:https?|ftp)://)" +
  // user:pass authentication
  "(?:\\S+(?::\\S*)?@)?" +
  "(?:" +
  // IP address exclusion
  // private & local networks
  "(?!(?:10|127)(?:\\.\\d{1,3}){3})" +
  "(?!(?:169\\.254|192\\.168)(?:\\.\\d{1,3}){2})" +
  "(?!172\\.(?:1[6-9]|2\\d|3[0-1])(?:\\.\\d{1,3}){2})" +
  // IP address dotted notation octets
  // excludes loopback network 0.0.0.0
  // excludes reserved space >= 224.0.0.0
  // excludes network & broacast addresses
  // (first & last IP address of each class)
  "(?:[1-9]\\d?|1\\d\\d|2[01]\\d|22[0-3])" +
  "(?:\\.(?:1?\\d{1,2}|2[0-4]\\d|25[0-5])){2}" +
  "(?:\\.(?:[1-9]\\d?|1\\d\\d|2[0-4]\\d|25[0-4]))" +
  "|" +
  // host name
  "(?:(?:[a-z\\u00a1-\\uffff0-9]-*)*[a-z\\u00a1-\\uffff0-9]+)" +
  // domain name
  "(?:\\.(?:[a-z\\u00a1-\\uffff0-9]-*)*[a-z\\u00a1-\\uffff0-9]+)*" +
  // TLD identifier
  "(?:\\.(?:[a-z\\u00a1-\\uffff]{2,}))" +
  ")" +
  // port number
  "(?::\\d{2,5})?" +
  // resource path
  "(?:/\\S*)?" +
  "$", "i"
);

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

var updateResults = function() {
  var $resultsContainer = $("#results-container");
  var $latestContainer = $("#latest-result-container");

  var source = $("#results-template").html();
  var template = Handlebars.compile(source);

  var latestSource = $("#latest-result-template").html();
  var latestTemplate = Handlebars.compile(latestSource);

  var url = $resultsContainer.data("url");
  var checkURL = $latestContainer.data("url");

  if ($resultsContainer.data("uninitialized")) {
    $resultsContainer.text("Loading…");
  }

  $.getJSON(url, function(data, status, xhr) {
    $.map(data, function(item, index) {
      item.URL = checkURL;

      if (item.Success) {
        item.Icon = "glyphicon-ok";
        item.TextStatus = "OK"
        item.CSSclass = "default"
        item.PanelCSSclass = "panel-success"
      } else if (item.status < 100) {
        item.Icon = "glyphicon-warning-sign";
        item.TextStatus = "Warning"
        item.CSSclass = "warning"
        item.PanelCSSclass = "panel-warning"
      } else {
        item.Icon = "glyphicon-fire";
        item.TextStatus = "Fail!"
        item.CSSclass = "danger"
        item.PanelCSSclass = "panel-danger"
      }
      return item
    });

    var html = template({"Results": data});
    $resultsContainer.html(html);
    dateRelativityApplicator($resultsContainer);

    var latestHtml = latestTemplate(data[0]);
    $latestContainer.html(latestHtml);
    dateRelativityApplicator($latestContainer);

    setTimeout(updateResults, 30*1000);
  });
}

$(function() {
  dateRelativityApplicator($("body"));

  var $checkForm = $("#check-form");
  if ($checkForm.length) {

    $checkForm.isHappy({
      fields: {
        "#url": {
          required: true,
          message: "You must enter a valid http or https URL. Private IP ranges are not allowed.",
          test: function (val) {
            return re_weburl.test(val);
          }
        },
      },
      unHappy: function() {
        $checkForm.find(".form-group").addClass("has-error").removeClass("has-success")
      },
      happy: function() {
        $checkForm.find(".form-group").addClass("has-success").removeClass("has-error")
      }
    });
  }

  $("#delete-check-form").on("submit", function(event){
    $button = $(this).find("button");
    if ($button.data("state") == "unarmed") {
      event.preventDefault()
      $button.data("state", "armed")
      $button.find("[data-role=label]").text("Are you sure?");
      setTimeout(function() {
        $button.data("state", "unarmed")
        $button.find("[data-role=label]").text("Destroy");
        $button.blur()
      }, 5000)
    }
  });

  var $resultsContainer = $("#results-container");
  if ($resultsContainer.length) {
    updateResults()
  }
});
