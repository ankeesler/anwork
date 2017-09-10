package com.marshmallow.anwork.event;

/**
 * This is an object that can hold onto {@link Event} instances.
 *
 * @author Andrew
 * Created Aug 31, 2017
 */
public interface EventLog {

  /**
   * Add a {@link Event} to this log.
   *
   * @param event The {@link Event} to add
   */
  public void add(Event event);

  /**
   * Clear all {@link Event}'s from the log.
   */
  public void clear();

  /**
   * Get the events that are currently in this log.
   *
   * The events should be returned in the order in which they were added. This
   * method should <b>never</b> return <code>null</code>.
   *
   * @return The events that are currently in this log
   */
  public Event[] getEvents();
}
