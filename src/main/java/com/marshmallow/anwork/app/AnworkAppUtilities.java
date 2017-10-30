package com.marshmallow.anwork.app;

import com.marshmallow.anwork.task.Task;

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
  public static final DateFormat DATE_FORMAT = new SimpleDateFormat("EEE MMM d HH:mm:ss");

  /**
   * Turn a {@link Task} into a human-readable one-line {@link String} that describes it.
   *
   * @param task The {@link Task} to turn into a {@link String}
   * @return A human-readable one-line {@link String} that describes the provided {@link Task}
   */
  public static String taskToString(Task task) {
    StringBuilder builder = new StringBuilder();
    builder.append(task.getName());
    builder.append(" (").append(task.getPriority()).append(")");
    builder.append(':');
    builder.append(" created ").append(DATE_FORMAT.format(task.getStartDate()));
    return builder.toString();
  }
}
