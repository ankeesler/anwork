package com.marshmallow.anwork.event;

import java.util.ArrayList;
import java.util.List;

/**
 * This is an {@link EventLog} where the events are stored in RAM.
 *
 * <p>
 * Created Aug 31, 2017
 * </p>
 *
 * @author Andrew
 */
public class RamEventLog implements EventLog {

  private List<Event> events = new ArrayList<Event>();

  @Override
  public void add(Event item) {
    events.add(item);
  }

  @Override
  public void clear() {
    events.clear();
  }

  @Override
  public Event[] getEvents() {
    return events.toArray(new Event[0]);
  }

}
