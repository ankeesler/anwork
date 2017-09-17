package com.marshmallow.anwork.journal;

import java.util.Date;

/**
 * This is an object that describes some activity that has happened.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public interface JournalEntry {

  /**
   * Get the title of this entry.
   *
   * @return The title of this entry
   */
  public String getTitle();

  /**
   * Get the description of this entry.
   *
   * @return The description of this entry
   */
  public String getDescription();

  /**
   * Get the date on which this activity happened.
   *
   * @return The date on which this activity happened
   */
  public Date getDate();
}
