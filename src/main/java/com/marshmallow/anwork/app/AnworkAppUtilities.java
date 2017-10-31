package com.marshmallow.anwork.app;

import com.marshmallow.anwork.journal.Journal;
import com.marshmallow.anwork.journal.JournalEntry;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;

import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;

/**
 * This class contains utilities shared throughout the ANWORK app.
 *
 * <p>
 * Created Oct 30, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkAppUtilities {

  private AnworkAppUtilities() { }

  /**
   * This is the common date format that should be used in displaying {@link Date} objects.
   */
  private static final DateFormat DATE_FORMAT = new SimpleDateFormat("EEE MMM d HH:mm:ss");

  /**
   * Turn a {@link Task} into a human-readable multi-line {@link String} that describes it.
   *
   * @param task The {@link Task} to turn into a {@link String}
   * @param manager The {@link TaskManager} associated with this {@link Task}
   * @param indent The indentation to apply to each line of the human-readable {@link String}
   *     returned by this function
   * @return A human-readable multi-line {@link String} that describes the provided {@link Task}
   */
  public static String makeTaskLongString(Task task, TaskManager manager, String indent) {
    StringBuilder builder = new StringBuilder();
    builder.append(indent)
           .append(task.getName())
           .append(' ')
           .append('(')
           .append(task.getId())
           .append(')')
           .append('\n');
    builder.append(indent)
           .append(indent)
           .append("created ")
           .append(DATE_FORMAT.format(task.getStartDate()))
           .append('\n');
    builder.append(indent)
           .append(indent)
           .append("priority ")
           .append(task.getPriority())
           .append('\n');
    builder.append(indent)
           .append(indent)
           .append('{')
           .append(getMostRecentJournalEntry(task, manager).getTitle())
           .append('}');
    return builder.toString();
  }

  private static JournalEntry getMostRecentJournalEntry(Task task, TaskManager manager) {
    // We assume that the provided Task does have a journal.
    Journal<?> journal = manager.getJournal(task.getName());
    JournalEntry[] entries = journal.getEntries();
    return entries[entries.length - 1];
  }

  /**
   * Turn a {@link JournalEntry} into a human-readable one-line {@link String} that describes it.
   *
   * @param entry The {@link JournalEntry} to turn into a {@link String}
   * @return A human-readable one-line {@link String} that describes the provided
   *     {@link JournalEntry}
   */
  public static String makeJournalEntryShortString(JournalEntry entry) {
    return String.format("[%s]: %s - %s",
                         DATE_FORMAT.format(entry.getDate()),
                         entry.getTitle(),
                         entry.getDescription());
  }
}
