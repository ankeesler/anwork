package com.marshmallow.anwork.smoketest;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

import org.junit.Test;

public class Smoketest {

  @Test
  public void runSanity() {
    assertTrue(true);
  }

  @Test
  public void runSanity0() {
    assertEquals(1, 1);
  }

  @Test
  public void runFailure() {
    assertEquals(1, 0);
  }
}
