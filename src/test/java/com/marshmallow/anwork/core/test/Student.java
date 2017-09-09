package com.marshmallow.anwork.core.test;

import com.marshmallow.anwork.core.Serializer;

/**
 * This is a dummy object to be used in tests.
 *
 * @author Andrew
 * @date Sep 8, 2017
 */
public class Student {

  private static class StudentSerializer implements Serializer<Student> {

    public static final StudentSerializer instance = new StudentSerializer();

    @Override
    public String marshall(Student t) {
      return String.format("Student:name=%s;id=%d;", t.getName(), t.getId());
    }

    @Override
    public Student unmarshall(String string) {
      StringBuffer buffer = new StringBuffer(string);
      if (!buffer.toString().startsWith("Student:")) {
        return null;
      }
      buffer.delete(0, "Student:".length());

      if (!buffer.toString().startsWith("name=")) {
        return null;
      }
      buffer.delete(0, "name=".length());
      int end = buffer.indexOf(";");
      if (end == -1) {
        return null;
      }
      String name = buffer.substring(0, end);
      buffer.delete(0, name.length() + 1);

      if (!buffer.toString().startsWith("id=")) {
        return null;
      }
      buffer.delete(0, "id=".length());
      end = buffer.indexOf(";");
      if (end == -1) {
        return null;
      }
      String id = buffer.substring(0, end);
      buffer.delete(0, id.length() + 1);

      if (buffer.length() != 0) {
        return null;
      }

      int idInt;
      try {
        idInt = Integer.parseInt(id);
      } catch (NumberFormatException nfe) {
        return null;
      }

      return new Student(name, idInt);
    }
  }

  public static Serializer<Student> serializer() {
    return StudentSerializer.instance;
  }

  private String name;
  private int id;

  public Student(String name, int id) {
    this.name = name;
    this.id = id;
  }

  public String getName() {
    return name;
  }

  public int getId() {
    return id;
  }

}
