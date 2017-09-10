package com.marshmallow.anwork.event.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import com.marshmallow.anwork.event.Event;
import com.marshmallow.anwork.event.EventLog;
import com.marshmallow.anwork.event.RamEventLog;

import java.util.Date;

import org.junit.Test;

public class EventLogTest {

  private static class TestEvent implements Event {

    private String type;
    private Date date;
    private String description;

    public TestEvent(String type, String description) {
      this.type = type;
      this.date = new Date();
      this.description = description;
    }

    @Override
    public String getType() {
      return type;
    }

    @Override
    public Date getDate() {
      return date;
    }

    @Override
    public String getDescription() {
      return description;
    }
  }

  private EventLog log = new RamEventLog();

  @Test
  public void basicAddTest() {
    Event[] events = log.getEvents();
    assertNotNull(events);
    assertEquals(0, events.length);

    TestEvent firstEvent = new TestEvent("type-a", "here is my first event");
    log.add(firstEvent);
    events = log.getEvents();
    assertNotNull(events);
    assertEquals(1, events.length);
    assertEquals(firstEvent, events[0]);

    TestEvent secondEvent = new TestEvent("type-b", "here is my second event");
    log.add(secondEvent);
    events = log.getEvents();
    assertNotNull(events);
    assertEquals(2, events.length);
    assertEquals(firstEvent, events[0]);
    assertEquals(secondEvent, events[1]);
  }

  @Test
  public void clearTest() {
    log.add(new TestEvent("type-1", "here is an event of type 1"));
    log.clear();
    assertEquals(0, log.getEvents().length);

    log.add(new TestEvent("type-1", "here is an event of type 1"));
    log.add(new TestEvent("type-2", "here is an event of type 2"));
    log.clear();
    assertEquals(0, log.getEvents().length);
  }
}
