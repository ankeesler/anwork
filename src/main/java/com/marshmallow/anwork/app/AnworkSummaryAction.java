package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Action;
import com.marshmallow.anwork.app.cli.ArgumentType;
import com.marshmallow.anwork.app.cli.ArgumentValues;
import com.marshmallow.anwork.journal.Journal;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskManagerActionType;
import com.marshmallow.anwork.task.TaskManagerJournalEntry;
import com.marshmallow.anwork.task.TaskState;

import java.util.Calendar;
import java.util.Date;

/**
 * This is a {@link Action} for the "anwork summary" CLI command.
 *
 * <p>
 * Created Nov 10, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkSummaryAction extends TaskManagerCliAction {
  @Override
  public void run(AnworkAppConfig config,
                  ArgumentValues flags,
                  ArgumentValues arguments,
                  TaskManager manager) {
    Long days = arguments.getValue("days", ArgumentType.NUMBER);
    Journal<TaskManagerJournalEntry> journal = manager.getJournal();
    TaskManagerJournalEntry[] entries = journal.getEntries();
    int recentEntriesCount = getRecentJournalEntriesCount(entries, days.intValue());
    printFinishedTasks(entries, recentEntriesCount);
  }

  // This method returns an integer describing how many of the journal entries happened in the last
  // "days" where "days" is the second parameter.
  private static int getRecentJournalEntriesCount(TaskManagerJournalEntry[] entries, int days) {
    Calendar calendar = Calendar.getInstance();
    calendar.add(Calendar.DATE, (days * -1));
    for (int i = entries.length - 1; i > -1; i--) {
      Date entryDate = entries[i].getDate();
      if (entryDate.compareTo(calendar.getTime()) < 0) {
        return entries.length - i;
      }
    }
    return entries.length;
  }

  // This method returns an array of Task's that have been finished
  private static void printFinishedTasks(TaskManagerJournalEntry[] entries,
                                         int recentEntriesCount) {
    for (int i = 0; i < recentEntriesCount; i++) {
      TaskManagerJournalEntry entry = entries[entries.length - i - 1];
      if (entry.getActionType().equals(TaskManagerActionType.SET_STATE)
          && entry.getDetail().equals(TaskState.FINISHED.name())) {
        System.out.printf("[%s]: Finished '%s'\n",
                          AnworkAppUtilities.DATE_FORMAT.format(entry.getDate()),
                          entry.getTask().getName());
        System.out.println("  took " + getTaskDuration(entry));
      }
    }
  }

  private static String getTaskDuration(TaskManagerJournalEntry finishedEntry) {
    Date startDate = finishedEntry.getTask().getStartDate();
    Date finishedDate = finishedEntry.getDate();

    long diffMs = finishedDate.getTime() - startDate.getTime();
    if (diffMs < 1000) {
      return String.format("%d milliseconds", diffMs);
    }

    long diffS = diffMs / 1000;
    if (diffS < 60) {
      return String.format("%d seconds", diffS);
    }

    long diffM = diffS / 60;
    if (diffM < 60) {
      return String.format("%d minutes", diffM);
    }

    long diffH = diffM / 60;
    if (diffH < 24) {
      return String.format("%h hours", diffH);
    }

    long diffD = diffH / 24;
    return String.format("%d days", diffD);
  }
}