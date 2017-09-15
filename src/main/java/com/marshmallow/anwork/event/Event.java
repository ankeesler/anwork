package com.marshmallow.anwork.event;

import java.util.Date;

/**
 * This is something that happened that you want to store in a {@link EventLog}.
 *
 * <p>
 * Created Aug 31, 2017
 * </p>
 *
 * @author Andrew
 */
public interface Event {

  /**
   * Get the type of this event.
   *
   * <p>
   * This is purposely loosely typed for academic purposes...
   * </p>
   *
   * @return The type of this event
   */
  public String getType();

  /**
   * Get the date at which this thing happened.
   *
   * @return The date at which this thing happened
   */
  public Date getDate();

  /**
   * Get the description of this event.
   *
   * @return The description of this event
   */
  public String getDescription();
}
