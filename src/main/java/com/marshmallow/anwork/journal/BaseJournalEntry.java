package com.marshmallow.anwork.journal;

import java.util.Date;

/**
 * This is a canonical {@link JournalEntry}. This object is meant to be subclassed.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public class BaseJournalEntry implements JournalEntry {

  private final String title;
  private final String description;
  private final Date date;

  /**
   * Create a canonical {@link JournalEntry} with a {@link #date} of now.
   *
   * @param title The title of this event
   * @param description The description of this event
   */
  public BaseJournalEntry(String title, String description) {
    this.title = title;
    this.description = description;
    this.date = new Date();
  }

  @Override
  public String getTitle() {
    return title;
  }

  @Override
  public String getDescription() {
    return description;
  }

  @Override
  public Date getDate() {
    return date;
  }

}
